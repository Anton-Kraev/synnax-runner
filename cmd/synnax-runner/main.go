package main

import (
	"log"

	"synnax-runner/internal/bot"
	"synnax-runner/internal/config"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Валидируем конфигурацию
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	log.Printf("🚀 Запуск Telegram бота...")

	// Создаем бота
	b, err := bot.New(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем бота
	b.Start()
}
