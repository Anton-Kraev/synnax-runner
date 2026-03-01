# Используем официальный образ Go
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/synnax-runner

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем Python и необходимые пакеты
RUN apk --no-cache add python3 py3-pip

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение
COPY --from=builder /app/main .

# Создаем директорию для скриптов
RUN mkdir -p scripts && chown -R appuser:appgroup /app

# Переключаемся на пользователя
USER appuser

# Открываем порт (если понадобится в будущем)
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
