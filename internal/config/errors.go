package config

import "errors"

var (
	// ErrMissingTelegramToken возникает, когда не указан токен Telegram бота
	ErrMissingTelegramToken = errors.New("TELEGRAM_BOT_TOKEN environment variable or -token flag is required")
)
