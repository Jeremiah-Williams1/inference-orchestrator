package routes

import (
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/response"
	"github.com/gin-gonic/gin"
)

func (s *routerImpl) GetHealth(c *gin.Context) {
	response.OK(c, "ok", nil)
}
