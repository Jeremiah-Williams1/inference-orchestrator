package job

import "log/slog"

type JobHandler struct {
	service *JobService
	log     *slog.Logger
}

func NewJobHandler(svc *JobService, log *slog.Logger) *JobHandler {
	return &JobHandler{service: svc, log: log}
}

// TODO: SubmitJob, GetJobResult, GetQueueDepth
