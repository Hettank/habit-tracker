package models

import "time"

type Habit struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Color       string
	Icon        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type HabitLog struct {
	ID          int64
	HabitID     int64
	CompletedAt time.Time
	CreatedAt   time.Time
}
