package routes

import (
	"net/http"

	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/api"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Job handlers live here.
// Each method maps to an operationId in api/openapi.yaml.
//
// TODO: add these endpoints to openapi.yaml first, run oapi-codegen,
// then implement the generated interface methods below.
//
// func (s *routerImpl) SubmitJob(c *gin.Context) {}
// func (s *routerImpl) GetJobResult(c *gin.Context, id string) {}
// func (s *routerImpl) GetQueueDepth(c *gin.Context, jobType string) {}

func (s *routerImpl) SubmitJob(c *gin.Context) {
	if s.job == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Service Unavaiilable", "code": "Server Error"})
		return
	}
	s.job.SubmitJob(c)
}

func (s *routerImpl) GetJobResult(c *gin.Context, id openapi_types.UUID) {
	if s.job == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Service Unavaiilable", "code": "Server Error"})
		return
	}
	s.job.GetJobResult(c, id)
}

func (s *routerImpl) GetQueueDepth(c *gin.Context, pType api.GetQueueDepthParamsType) {
	if s.job == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Service Unavaiilable", "code": "Server Error"})
		return
	}
	s.job.GetQueueDepth(c, pType)
}
