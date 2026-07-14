package dto

type DashboardResponse struct {
	TotalHabits          int64 `json:"total_habits"`
	CompletedToday       int64 `json:"completed_today"`
	PendingToday         int64 `json:"pending_today"`
	CompletionPercentage int64 `json:"completion_percentage"`
}
