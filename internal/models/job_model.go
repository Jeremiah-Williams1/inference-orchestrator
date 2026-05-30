package models

import "time"

// Status represents where a job is in its lifecycle.
// A job moves in one direction: Queued -> Processing -> Completed (or Failed).
type Status string

const (
	StatusQueued     Status = "queued"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Type tells the system which worker pool should handle this job.
// Each model type has its own Redis queue and its own pool of workers.
type Type string

const (
	TypeClassification Type = "classification"
	TypeRegression     Type = "regression" // Phase 2
)

// Job is the unit of work that flows through the system.
// The API creates it. Redis carries it. The worker processes it.
// The result lives back in Redis until the caller retrieves it.
type Job struct {
	ID        string                 `json:"id"`
	Type      Type                   `json:"type"`
	Status    Status                 `json:"status"`
	Input     map[string]interface{} `json:"input"`
	Result    interface{}            `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
