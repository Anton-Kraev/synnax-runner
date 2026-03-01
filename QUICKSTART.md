# Быстрый старт

## 1. Создание Telegram бота

1. Откройте Telegram и найдите @BotFather
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Сохраните полученный токен

## 2. Запуск бота

### Локальный запуск

```bash
# Установите зависимости
go mod tidy

# Способ 1: Через переменную окружения
export TELEGRAM_BOT_TOKEN="ваш_токен_здесь"
go run ./cmd/synnax-runner

# Способ 2: Через CLI флаг
go run ./cmd/synnax-runner -token "ваш_токен_здесь"
# или
go run ./cmd/synnax-runner -t "ваш_токен_здесь"
```

### Запуск с виртуальным окружением Python (venv)

Для выполнения `.py` и `.ipynb` скриптов бот вызывает `python3` и `python3 -m jupyter`. Чтобы использовать интерпретатор из venv (с установленным Jupyter), можно сделать так:

**Вариант А: запуск из активированного venv**

```bash
# Создайте venv и установите Jupyter (один раз)
python3 -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install jupyter nbconvert

# В этом же терминале запустите бота (будет использоваться python из venv)
export TELEGRAM_BOT_TOKEN="ваш_токен_здесь"
make run
# или: go run ./cmd/synnax-runner
```

**Вариант Б: указать путь к Python (systemd, cron, без активации)**

```bash
# Укажите полный путь к python в venv
export TELEGRAM_BOT_TOKEN="ваш_токен_здесь"
export PYTHON_PATH=".venv/bin/python3"   # или абсолютный путь: /path/to/project/.venv/bin/python3
make run
```

Через `.env` (если приложение его подхватывает) или в systemd unit:
```
Environment="PYTHON_PATH=/path/to/project/.venv/bin/python3"
```

### Запуск через Makefile

```bash
# Установите зависимости
make deps

# Запустите бота (для .ipynb предварительно активируйте venv или задайте PYTHON_PATH)
make run
```

### Запуск через Docker

```bash
# Создайте файл .env с токеном
echo "TELEGRAM_BOT_TOKEN=ваш_токен_здесь" > .env

# Запустите через docker-compose
docker-compose up -d
```

## 3. Использование бота

1. Найдите вашего бота в Telegram
2. Отправьте `/start`
3. Отправьте `/schedule`
4. Настройте расписание выполнения (например: `0 6 * * *` для ежедневного выполнения в 6:00 UTC)
5. Установите таймаут выполнения (например: `300` для 5 минут)
6. Загрузите Python файл (.py) или Jupyter Notebook (.ipynb)
7. Бот автоматически запланирует выполнение согласно вашим настройкам

## 4. Тестирование

```bash
# Запустите тесты
go test -v

# Или через Makefile
make test
```

## 5. Остановка

### Локальный запуск
Нажмите `Ctrl+C`

### Docker
```bash
docker-compose down
```
