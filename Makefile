.PHONY: build run test clean deps

# Переменные
BINARY_NAME=synnax-runner
BUILD_DIR=build

# Сборка проекта
build:
	@echo "🔨 Сборка проекта..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/synnax-runner

# Запуск проекта
run:
	@echo "🚀 Запуск бота..."
	go run ./cmd/synnax-runner

# Установка зависимостей
deps:
	@echo "📦 Установка зависимостей..."
	go mod tidy
	go mod download

# Тестирование
test:
	@echo "🧪 Запуск тестов..."
	go test ./...

# Очистка
clean:
	@echo "🧹 Очистка..."
	rm -rf $(BUILD_DIR)
	rm -rf scripts/
	go clean

# Проверка кода
lint:
	@echo "🔍 Проверка кода..."
	golangci-lint run

# Форматирование кода
fmt:
	@echo "🎨 Форматирование кода..."
	go fmt ./...

# Проверка зависимостей на уязвимости
security:
	@echo "🔒 Проверка безопасности..."
	go list -json -deps ./... | nancy sleuth

# Создание директории для скриптов
setup:
	@echo "📁 Создание директорий..."
	mkdir -p scripts

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build     - Сборка проекта"
	@echo "  run       - Запуск бота"
	@echo "  deps      - Установка зависимостей"
	@echo "  test      - Запуск тестов"
	@echo "  clean     - Очистка проекта"
	@echo "  lint      - Проверка кода"
	@echo "  fmt       - Форматирование кода"
	@echo "  security  - Проверка безопасности"
	@echo "  setup     - Создание директорий"
	@echo "  help      - Показать эту справку"
