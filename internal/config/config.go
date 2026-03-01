package config

import (
	"flag"
	"os"
)

// Config содержит все настройки приложения
type Config struct {
	TelegramToken string
}

// Load загружает конфигурацию из переменных окружения и CLI флагов
func Load() *Config {
	config := &Config{}

	// Определяем CLI флаги
	flag.StringVar(&config.TelegramToken, "token", "", "Telegram Bot Token")
	flag.StringVar(&config.TelegramToken, "t", "", "Telegram Bot Token (short)")
	flag.Parse()

	// Если токен не передан через флаг, берем из переменной окружения
	if config.TelegramToken == "" {
		config.TelegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	}

	return config
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.TelegramToken == "" {
		return ErrMissingTelegramToken
	}
	return nil
}
