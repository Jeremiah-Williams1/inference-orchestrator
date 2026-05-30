// Package routes wires dependencies and registers all HTTP routes.
// This mirrors the pattern from your previous project:
//   - Router struct holds all dependencies
//   - New() is the constructor called from main
//   - Routes() returns the configured gin.Engine
//   - Each resource has its own file (health.go, job.go etc)
//     but all implement methods on the same routerImpl struct,
//     which satisfies the generated api.ServerInterface.
package routes

import (
	"log/slog"
	"os"

	"github.com/Jeremiah-Williams1/inference-orchestrator/config"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/docs"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/job"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/middleware"
	redis "github.com/Jeremiah-Williams1/inference-orchestrator/pkg/redisclient"
	"github.com/gin-gonic/gin"
)

// Router holds every dependency the application needs.
// Add new dependencies here as the project grows.
type Router struct {
	cfg   *config.Config
	log   *slog.Logger
	redis *redis.Client
}

func New(cfg *config.Config, log *slog.Logger, redisClient *redis.Client) *Router {
	return &Router{cfg: cfg, log: log, redis: redisClient}
}

type routerImpl struct {
	job *job.JobHandler
}

func (s *Router) Routes() *gin.Engine {
	if s.cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	impl := &routerImpl{}

	// Health lives outside /api/v1 — never versioned, always stable.
	r.GET("/health", impl.GetHealth)

	spec, err := os.ReadFile("api.yaml")
	if err != nil {
		s.log.Warn("could not read api.yaml — /docs/spec will be unavailable", "err", err)
	}
	docs := docs.NewDocsHandler(spec)
	r.GET("/docs", docs.UI)
	r.GET("/docs/spec", docs.Spec)

	// TODO: once you run oapi-codegen and internal/api/gen.go exists,
	// replace the manual route above with:
	//
	// v1 := r.Group("/api/v1")
	// api.RegisterHandlersWithOptions(v1, impl, api.GinServerOptions{})
	//
	// After that, all routes are registered automatically from the generated code.
	// You only implement the methods — never touch route registration manually.

	return r
}
