package dto

import "time"

type CreateHabitRequest struct {
	Title       string `json:"title" validate:"required,max=25"`
	Description string `json:"description" validate:"max=250"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
}

type HabitResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateHabitRequest struct {
	Title       string `json:"title" validate:"required,max=25"`
	Description string `json:"description" validate:"max=250"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
}

type HabitLogResponse struct {
	ID          int64     `json:"id"`
	HabitID     int64     `json:"habit_id"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type HabitStreakResponse struct {
	CurrentStreak int64 `json:"current_streak"`
}

