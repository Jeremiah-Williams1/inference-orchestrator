package job

import (
	"log/slog"

	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/api"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/models"
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/response"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type SubmitJobRequest struct {
	Type  models.Type            `json:"type"`
	Input map[string]interface{} `json:"input" binding:"required"`
}

type JobHandler struct {
	service *JobService
	log     *slog.Logger
}

func NewJobHandler(svc *JobService, log *slog.Logger) *JobHandler {
	return &JobHandler{service: svc, log: log}
}

// TODO: SubmitJob, GetJobResult, GetQueueDepth
func (h *JobHandler) SubmitJob(c *gin.Context) {
	var req SubmitJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid request body", "error", err)
		response.ValidationErr(c, err)
		return
	}

	job, err := h.service.CreateJob(c.Request.Context(), req.Type, req.Input)
	if err != nil {
		h.log.Error("failed to create job", "error", err)
		response.Err(c, err)
		return
	}

	response.Accepted(c, "job queued", job)
}

func (h *JobHandler) GetJobResult(c *gin.Context, id openapi_types.UUID) {

	job, err := h.service.GetJob(c.Request.Context(), id.String())
	if err != nil {
		h.log.Error("failed to Get job", "error", err)
		response.Err(c, err)
		return
	}

	response.OK(c, "Job Gotten", job)

}

func (h *JobHandler) GetQueueDepth(c *gin.Context, pType api.GetQueueDepthParamsType) {
	depth, err := h.service.GetQueueDepth(c.Request.Context(), models.Type(pType))
	if err != nil {
		h.log.Error("failed to Get queue depth", "error", err)
		response.Err(c, err)
		return
	}

	response.OK(c, "queue depth", gin.H{"depth": depth, "type": pType})

}
