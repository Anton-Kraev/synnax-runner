package scheduler

import (
	"testing"
)

func TestValidateSchedule(t *testing.T) {
	tests := []struct {
		name    string
		schedule string
		wantErr bool
	}{
		{
			name:     "valid daily schedule",
			schedule: "0 6 * * *",
			wantErr:  false,
		},
		{
			name:     "valid weekly schedule",
			schedule: "0 9 * * 1-5",
			wantErr:  false,
		},
		{
			name:     "valid monthly schedule",
			schedule: "0 12 1 * *",
			wantErr:  false,
		},
		{
			name:     "invalid schedule",
			schedule: "invalid",
			wantErr:  true,
		},
		{
			name:     "empty schedule",
			schedule: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSchedule(tt.schedule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetScheduleDescription(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		want     string
	}{
		{
			name:     "daily schedule",
			schedule: "0 6 * * *",
			want:     "Каждый день в 6:00 UTC",
		},
		{
			name:     "weekly schedule",
			schedule: "0 9 * * 1-5",
			want:     "По будням в 9:00 UTC",
		},
		{
			name:     "custom schedule",
			schedule: "0 15 * * 2",
			want:     "По расписанию: 0 15 * * 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetScheduleDescription(tt.schedule)
			if got != tt.want {
				t.Errorf("GetScheduleDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_AddJob(t *testing.T) {
	scheduler := New()
	scheduler.Start()
	defer scheduler.Stop()

	// Тестируем добавление задачи
	jobID, err := scheduler.AddJob("* * * * *", func() {
		// Пустая функция для тестирования
	})

	if err != nil {
		t.Errorf("AddJob() error = %v", err)
	}

	if jobID == 0 {
		t.Error("AddJob() returned zero job ID")
	}
}
