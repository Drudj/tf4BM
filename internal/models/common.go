package models

import (
	"time"
)

// APIResponse представляет базовую структуру ответа API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError представляет ошибку API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error реализует интерфейс error
func (e *APIError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}

// PaginationMeta содержит метаданные пагинации
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

// PaginatedResponse представляет ответ с пагинацией
type PaginatedResponse struct {
	Data []interface{}  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// TaskStatus представляет статус задачи
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// Task представляет задачу в API
type Task struct {
	UUID        string     `json:"uuid"`
	Type        string     `json:"type"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"`
	Message     string     `json:"message,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// IsCompleted проверяет, завершена ли задача
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// IsFailed проверяет, провалилась ли задача
func (t *Task) IsFailed() bool {
	return t.Status == TaskStatusFailed || t.Status == TaskStatusCancelled
}

// IsInProgress проверяет, выполняется ли задача
func (t *Task) IsInProgress() bool {
	return t.Status == TaskStatusPending || t.Status == TaskStatusRunning
}
