package job

import (
	"log/slog"

	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/queue"
)

type JobService struct {
	queue queue.Queue
	log   *slog.Logger
}

func NewJobService(q queue.Queue, log *slog.Logger) *JobService {
	return &JobService{queue: q, log: log}
}

// TODO: CreateJob, GetJob, GetQueueDepth
