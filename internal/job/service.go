package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/models"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/queue"
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/apperror"
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/metrics"
	"github.com/google/uuid"
)

type JobService struct {
	queue   queue.Queue
	log     *slog.Logger
	metrics *metrics.Metrics
}

func NewJobService(q queue.Queue, log *slog.Logger, metrics *metrics.Metrics) *JobService {
	return &JobService{queue: q, log: log, metrics: metrics}
}

func (j JobService) CreateJob(ctx context.Context, jobType models.Type, input map[string]interface{}) (*models.Job, error) {
	id := uuid.NewString()

	job := models.Job{
		ID:        id,
		Type:      jobType,
		Input:     input,
		Status:    models.StatusQueued,
		CreatedAt: time.Now(),
	}

	if jobType != models.TypeClassification && jobType != models.TypeRegression {
		return nil, apperror.BadRequest("invalid job type")
	}

	err := j.queue.SetResult(ctx, &job)
	if err != nil {
		j.log.Error("Failed to set state for job", "job_id", job.ID, "error", err)
		return nil, fmt.Errorf("failed to Set state for  job %s: %w", job.ID, err)
	}

	err = j.queue.Enqueue(ctx, &job)
	if err != nil {
		j.log.Error("Failed to enqueue job", "job_id", job.ID, "error", err)
		return nil, fmt.Errorf("failed to enqueue job %s: %w", job.ID, err)
	}

	j.metrics.JobsSubmitted.WithLabelValues(string(jobType)).Inc()
	return &job, nil

}

func (j JobService) GetJob(ctx context.Context, jobID string) (*models.Job, error) {
	val, err := j.queue.GetResult(ctx, jobID)
	if err != nil {
		j.log.Error("Failed to Get result", "job_id", jobID, "error", err)
		return nil, fmt.Errorf("Couldn't get result for job %s, Error: %w", jobID, err)
	}

	if val == nil {
		return nil, apperror.NotFound("job not found")
	}

	return val, nil

}

func (j JobService) GetQueueDepth(ctx context.Context, jobType models.Type) (int64, error) {
	depth, err := j.queue.Depth(ctx, jobType)
	if err != nil {
		j.log.Error("Error getting length for", "job_type", jobType, "error", err)
		return 0, fmt.Errorf("Error gettin the length %w", err)
	}
	j.metrics.QueueDepth.WithLabelValues(string(jobType)).Set(float64(depth))
	return depth, nil
}
