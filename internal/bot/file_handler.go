package bot

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/telebot.v3"

	"synnax-runner/internal/scheduler"
	"synnax-runner/pkg/utils"
)

// handleFileUpload обрабатывает загрузку файлов
func (b *Bot) handleFileUpload(c telebot.Context) error {
	doc := c.Message().Document
	if doc == nil {
		return c.Send("Пожалуйста, отправьте файл, а не текст.")
	}

	ext := strings.ToLower(filepath.Ext(doc.FileName))
	if ext != ".py" && ext != ".ipynb" {
		return c.Send("Пожалуйста, отправьте Python файл (.py) или Jupyter Notebook (.ipynb)")
	}

	// Создаем директорию для файлов
	scriptsDir := "scripts"
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return c.Send("Ошибка при создании директории для файла")
	}

	// Скачиваем файл
	filePath := filepath.Join(scriptsDir, doc.FileName)

	// Создаем файл локально
	outputFile, err := os.Create(filePath)
	if err != nil {
		return c.Send("Ошибка при создании файла")
	}
	defer outputFile.Close()

	// Скачиваем содержимое файла
	reader, err := b.bot.File(&doc.File)
	if err != nil {
		return c.Send("Ошибка при получении файла")
	}

	// Копируем содержимое
	_, err = io.Copy(outputFile, reader)
	if err != nil {
		return c.Send("Ошибка при сохранении файла")
	}

	// Останавливаем предыдущую задачу
	b.stopActiveJob()

	// Сохраняем путь к файлу
	b.mu.Lock()
	b.userFile = filePath
	b.state = StateIdle
	b.mu.Unlock()

	// Планируем новую задачу
	b.scheduleJob(filePath)

	// Формируем сообщение об успехе
	scheduleDesc := scheduler.GetScheduleDescription(b.userSettings.Schedule)
	timeoutDesc := fmt.Sprintf("%d секунд", b.userSettings.Timeout)

	successMsg := fmt.Sprintf("✅ Файл `%s` успешно загружен!\n\n📅 **Расписание:** %s\n⏱️ **Таймаут:** %s\n\n🔄 Предыдущие задачи остановлены.",
		doc.FileName, scheduleDesc, timeoutDesc)

	return c.Send(successMsg)
}

// stopActiveJob останавливает активную задачу
func (b *Bot) stopActiveJob() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.activeJob != 0 {
		b.scheduler.RemoveJob(b.activeJob)
		b.activeJob = 0
	}
}

// scheduleJob планирует новую задачу
func (b *Bot) scheduleJob(filePath string) {
	// Используем расписание пользователя или значение по умолчанию
	schedule := b.userSettings.Schedule
	if schedule == "" {
		schedule = "0 6 * * *" // По умолчанию каждый день в 6:00 UTC
	}

	// Планируем выполнение
	jobID, err := b.scheduler.AddJob(schedule, func() {
		b.executeScript(filePath)
	})

	if err != nil {
		log.Printf("Ошибка при планировании задачи: %v", err)
		return
	}

	b.mu.Lock()
	b.activeJob = jobID
	b.mu.Unlock()
}

// executeScript выполняет скрипт (.py или .ipynb)
func (b *Bot) executeScript(filePath string) {
	ext := strings.ToLower(filepath.Ext(filePath))
	scriptType := "Python скрипта"
	if ext == ".ipynb" {
		scriptType = "Jupyter Notebook"
	}
	b.sendNotification("🚀 Начинаю выполнение " + scriptType + "...")

	result := b.executor.Execute(filePath, b.userSettings.Timeout)
	b.sendJobResult(result)
}

// sendJobResult отправляет результат выполнения
func (b *Bot) sendJobResult(result JobResult) {
	if result.Success {
		// Отправляем основное сообщение об успехе
		b.sendNotification("✅ Скрипт выполнен успешно!")

		// Отправляем вывод скрипта отдельным сообщением
		if result.Output != "" {
			outputMsg := "📋 Вывод скрипта:\n```\n" + utils.EscapeMarkdown(result.Output) + "\n```"
			b.sendNotification(outputMsg)
		}
	} else {
		// Отправляем основное сообщение об ошибке
		b.sendNotification("❌ Ошибка при выполнении скрипта!")

		// Отправляем ошибку отдельным сообщением
		if result.Error != "" {
			errorMsg := "🚨 Ошибка:\n```\n" + utils.EscapeMarkdown(result.Error) + "\n```"
			b.sendNotification(errorMsg)
		}

		// Отправляем вывод скрипта отдельным сообщением
		if result.Output != "" {
			outputMsg := "📋 Вывод скрипта:\n```\n" + utils.EscapeMarkdown(result.Output) + "\n```"
			b.sendNotification(outputMsg)
		}
	}
}

// sendNotification отправляет уведомление авторизованному пользователю
func (b *Bot) sendNotification(text string) {
	if b.authorizedChat == 0 {
		log.Printf("Нет авторизованного чата для отправки уведомления")
		return
	}

	// Telegram имеет лимит на длину сообщения (4096 символов)
	const maxMessageLength = 4000 // Оставляем небольшой запас

	if len(text) <= maxMessageLength {
		_, err := b.bot.Send(&telebot.Chat{ID: b.authorizedChat}, text, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
		if err != nil {
			log.Printf("Ошибка отправки уведомления: %v", err)
		}
	} else {
		// Разбиваем длинное сообщение на части
		parts := utils.SplitMessage(text, maxMessageLength)
		for i, part := range parts {
			_, err := b.bot.Send(&telebot.Chat{ID: b.authorizedChat}, part, &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
			if err != nil {
				log.Printf("Ошибка отправки части уведомления %d: %v", i+1, err)
			}

			// Небольшая задержка между сообщениями
			time.Sleep(100 * time.Millisecond)
		}
	}
}
