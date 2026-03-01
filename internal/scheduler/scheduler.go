package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler управляет планированием задач
type Scheduler struct {
	cron *cron.Cron
}

// New создает новый планировщик
func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithLocation(time.UTC)),
	}
}

// Start запускает планировщик
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop останавливает планировщик
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// AddJob добавляет новую задачу
func (s *Scheduler) AddJob(schedule string, job func()) (cron.EntryID, error) {
	return s.cron.AddFunc(schedule, job)
}

// RemoveJob удаляет задачу
func (s *Scheduler) RemoveJob(jobID cron.EntryID) {
	s.cron.Remove(jobID)
}

// ValidateSchedule проверяет корректность cron выражения
func ValidateSchedule(schedule string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(schedule)
	return err
}

// GetScheduleDescription возвращает описание расписания
func GetScheduleDescription(schedule string) string {
	// Простые описания для популярных расписаний
	switch schedule {
	case "0 6 * * *":
		return "Каждый день в 6:00 UTC"
	case "0 9 * * 1-5":
		return "По будням в 9:00 UTC"
	case "*/30 * * * *":
		return "Каждые 30 минут"
	case "0 12 1 * *":
		return "1-го числа каждого месяца в 12:00 UTC"
	case "0 0 * * 0":
		return "Каждое воскресенье в полночь"
	case "0 */2 * * *":
		return "Каждые 2 часа"
	case "0 0 1 * *":
		return "1-го числа каждого месяца в полночь"
	default:
		return "По расписанию: " + schedule
	}
}
