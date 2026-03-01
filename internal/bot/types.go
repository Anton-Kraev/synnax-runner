package bot

// UserSettings содержит настройки пользователя
type UserSettings struct {
	Schedule string // cron schedule
	Timeout  int    // timeout in seconds
}

// JobResult содержит результат выполнения задачи
type JobResult struct {
	Success bool
	Output  string
	Error   string
}

// BotState представляет состояние бота
type BotState string

const (
	StateWaitingSchedule BotState = "waiting_schedule"
	StateWaitingTimeout  BotState = "waiting_timeout"
	StateWaitingFile     BotState = "waiting_file"
	StateIdle            BotState = "idle"
)
