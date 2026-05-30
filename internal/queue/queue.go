// Package queue defines the contract between the service layer and the queue backend.
//
// The Queue interface is what the service talks to.
// The Redis implementation is what actually runs.
// The service never imports Redis directly — it imports this interface.
//
// Why does this matter?
// When you write tests for the service, you pass in a fake queue that stores
// jobs in memory. No Redis needed, no Docker needed, tests run instantly.
// The service cannot tell the difference because it only sees the interface.
//
// This is the same pattern your previous codebase enforced with interfaces on routes.
// Same idea, different layer.
package queue

import (
	"context"

	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/models"
)

// Queue is the contract every queue implementation must satisfy.
// Right now there is one implementation: Redis.
// If you ever want SQS, Kafka, or an in-memory queue for tests — implement this interface.
type Queue interface {
	// Enqueue pushes a job onto the queue for its type.
	// classification jobs go to one queue, regression jobs go to another.
	Enqueue(ctx context.Context, j *models.Job) error

	// Dequeue blocks until a job is available and returns it.
	// Workers call this in their main processing loop.
	// Returns nil, nil when the timeout is reached and no job was available — this is normal.
	Dequeue(ctx context.Context, jobType models.Type) (*models.Job, error)

	// SetResult stores a job's current state so the handler can retrieve it.
	// Called once when the job is created (status: queued).
	// Called again by the worker when processing completes (status: completed or failed).
	SetResult(ctx context.Context, j *models.Job) error

	// GetResult retrieves a job's current state by ID.
	// Returns nil, nil when the job ID does not exist.
	GetResult(ctx context.Context, jobID string) (*models.Job, error)

	// Depth returns how many jobs are currently waiting in the queue for a given type.
	// KEDA polls this to decide when to scale workers up or down.
	Depth(ctx context.Context, jobType models.Type) (int64, error)
}
