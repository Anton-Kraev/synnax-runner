package config

import (
	"os"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "test_token",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				TelegramToken: tt.token,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad_FromEnvironment(t *testing.T) {
	// Сохраняем оригинальное значение
	originalToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	defer os.Setenv("TELEGRAM_BOT_TOKEN", originalToken)

	// Устанавливаем тестовое значение
	testToken := "test_token_from_env"
	os.Setenv("TELEGRAM_BOT_TOKEN", testToken)

	// Загружаем конфигурацию
	cfg := Load()

	if cfg.TelegramToken != testToken {
		t.Errorf("Expected token %s, got %s", testToken, cfg.TelegramToken)
	}
}
