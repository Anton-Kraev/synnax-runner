package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v3"

	"synnax-runner/internal/scheduler"
)

// Bot представляет Telegram бота
type Bot struct {
	bot            *telebot.Bot
	scheduler      *scheduler.Scheduler
	executor       *Executor
	authorizedUser int64 // ID авторизованного пользователя
	authorizedChat int64 // Chat ID авторизованного пользователя
	state          BotState
	userFile       string
	userSettings   UserSettings
	mu             sync.RWMutex
	activeJob      cron.EntryID
}

// New создает нового бота
func New(token string) (*Bot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 1 * time.Minute},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:            bot,
		scheduler:      scheduler.New(),
		executor:       NewExecutor(),
		authorizedUser: 0,
		authorizedChat: 0,
		state:          StateIdle,
		userFile:       "",
		userSettings:   UserSettings{},
		activeJob:      0,
	}, nil
}

// Start запускает бота
func (b *Bot) Start() {
	b.scheduler.Start()
	log.Printf("Bot started: @%s", b.bot.Me.Username)

	// Регистрируем обработчики команд
	b.bot.Handle("/start", b.handleStart)
	b.bot.Handle("/schedule", b.handleSchedule)
	b.bot.Handle("/help", b.handleHelp)
	b.bot.Handle(telebot.OnText, b.handleText)
	b.bot.Handle(telebot.OnDocument, b.handleDocument)

	// Запускаем бота
	b.bot.Start()
}

// Stop останавливает бота
func (b *Bot) Stop() {
	b.scheduler.Stop()
}

// checkAuthorization проверяет авторизацию пользователя
func (b *Bot) checkAuthorization(c telebot.Context) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	userID := c.Sender().ID
	chatID := c.Chat().ID

	// Если это первый пользователь, авторизуем его
	if b.authorizedUser == 0 {
		b.authorizedUser = userID
		b.authorizedChat = chatID
		log.Printf("Авторизован пользователь: %d в чате: %d", userID, chatID)
		return true
	}

	// Проверяем, является ли пользователь авторизованным
	if b.authorizedUser != userID {
		return false
	}

	return true
}

// handleStart обрабатывает команду /start
func (b *Bot) handleStart(c telebot.Context) error {
	if !b.checkAuthorization(c) {
		return c.Send("⛔ Доступ запрещен. Бот уже используется другим пользователем.")
	}

	return c.Send("Привет! Я бот для планирования Python скриптов.\n\nИспользуйте /schedule для планирования нового скрипта.")
}

// handleSchedule обрабатывает команду /schedule
func (b *Bot) handleSchedule(c telebot.Context) error {
	if !b.checkAuthorization(c) {
		return c.Send("⛔ Доступ запрещен. Бот уже используется другим пользователем.")
	}

	b.mu.Lock()
	b.state = StateWaitingSchedule
	b.mu.Unlock()

	return c.Send("Настройте расписание выполнения скрипта.\n\nПримеры:\n• `0 6 * * *` - каждый день в 6:00 UTC\n• `0 9 * * 1-5` - по будням в 9:00 UTC\n• `*/30 * * * *` - каждые 30 минут\n• `0 12 1 * *` - 1-го числа каждого месяца в 12:00 UTC\n\nОтправьте расписание в формате cron:")
}

// handleHelp обрабатывает команду /help
func (b *Bot) handleHelp(c telebot.Context) error {
	if !b.checkAuthorization(c) {
		return c.Send("⛔ Доступ запрещен. Бот уже используется другим пользователем.")
	}

	return b.sendHelpMessage(c)
}

// handleText обрабатывает текстовые сообщения
func (b *Bot) handleText(c telebot.Context) error {
	if !b.checkAuthorization(c) {
		return c.Send("⛔ Доступ запрещен. Бот уже используется другим пользователем.")
	}

	b.mu.Lock()
	state := b.state
	b.mu.Unlock()

	switch state {
	case StateWaitingSchedule:
		return b.handleScheduleInput(c)
	case StateWaitingTimeout:
		return b.handleTimeoutInput(c)
	default:
		return c.Send("Используйте /schedule для планирования нового скрипта или /help для справки.")
	}
}

// handleDocument обрабатывает загрузку файлов
func (b *Bot) handleDocument(c telebot.Context) error {
	if !b.checkAuthorization(c) {
		return c.Send("⛔ Доступ запрещен. Бот уже используется другим пользователем.")
	}

	b.mu.Lock()
	state := b.state
	b.mu.Unlock()

	if state == StateWaitingFile {
		return b.handleFileUpload(c)
	}

	return c.Send("Пожалуйста, сначала настройте расписание и таймаут с помощью команды /schedule.")
}

// handleScheduleInput обрабатывает ввод расписания
func (b *Bot) handleScheduleInput(c telebot.Context) error {
	schedule := strings.TrimSpace(c.Text())

	// Проверяем, что это не команда
	if strings.HasPrefix(schedule, "/") {
		return c.Send("Пожалуйста, отправьте расписание в формате cron, а не команду.")
	}

	// Валидируем cron выражение
	if err := scheduler.ValidateSchedule(schedule); err != nil {
		return c.Send(fmt.Sprintf("❌ Неверный формат расписания: %v\n\nПопробуйте еще раз или используйте /help для примеров.", err))
	}

	// Сохраняем расписание
	b.mu.Lock()
	b.userSettings.Schedule = schedule
	b.state = StateWaitingTimeout
	b.mu.Unlock()

	return c.Send(fmt.Sprintf("✅ Расписание установлено: `%s`\n\nТеперь установите таймаут выполнения скрипта (в секундах):\n\nПримеры:\n• `300` - 5 минут\n• `600` - 10 минут\n• `1800` - 30 минут\n• `3600` - 1 час", schedule))
}

// handleTimeoutInput обрабатывает ввод таймаута
func (b *Bot) handleTimeoutInput(c telebot.Context) error {
	timeoutStr := strings.TrimSpace(c.Text())

	// Проверяем, что это не команда
	if strings.HasPrefix(timeoutStr, "/") {
		return c.Send("Пожалуйста, отправьте таймаут в секундах, а не команду.")
	}

	// Парсим таймаут
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout <= 0 {
		return c.Send("❌ Неверный таймаут. Пожалуйста, отправьте положительное число в секундах.")
	}

	// Проверяем разумные пределы (от 30 секунд до 1 часа)
	if timeout < 30 || timeout > 86400 {
		return c.Send("❌ Таймаут должен быть от 30 секунд до 1 часа (3600 секунд).")
	}

	// Сохраняем таймаут
	b.mu.Lock()
	b.userSettings.Timeout = timeout
	b.state = StateWaitingFile
	b.mu.Unlock()

	return c.Send(fmt.Sprintf("✅ Таймаут установлен: %d секунд\n\nТеперь отправьте Python файл (.py) или Jupyter Notebook (.ipynb)", timeout))
}

// sendHelpMessage отправляет справку
func (b *Bot) sendHelpMessage(c telebot.Context) error {
	helpText := "🤖 **Справка по боту**\n\n" +
		"**Команды:**\n" +
		"• `/start` - Начало работы\n" +
		"• `/schedule` - Планирование нового скрипта\n" +
		"• `/help` - Показать эту справку\n\n" +
		"**Процесс планирования:**\n" +
		"1. `/schedule` - Настройка расписания\n" +
		"2. Установка таймаута выполнения\n" +
		"3. Загрузка Python файла\n\n" +
		"**Примеры расписаний (cron):**\n" +
		"• `0 6 * * *` - каждый день в 6:00 UTC\n" +
		"• `0 9 * * 1-5` - по будням в 9:00 UTC\n" +
		"• `*/30 * * * *` - каждые 30 минут\n" +
		"• `0 12 1 * *` - 1-го числа каждого месяца в 12:00 UTC\n" +
		"• `0 0 * * 0` - каждое воскресенье в полночь\n\n" +
		"**Формат cron:**\n" +
		"`минута час день_месяца месяц день_недели`\n\n" +
		"**Таймауты:**\n" +
		"• От 30 секунд до 1 часа\n" +
		"• Рекомендуется: 300-1800 секунд\n\n" +
		"**Примечания:**\n" +
		"• Каждый новый скрипт отменяет предыдущие\n" +
		"• Время указано в UTC\n" +
		"• Поддерживаются .py и .ipynb (Jupyter Notebook) файлы\n" +
		"• Бот работает только с одним пользователем"

	return c.Send(helpText)
}
