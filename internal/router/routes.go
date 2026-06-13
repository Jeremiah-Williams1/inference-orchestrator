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
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/api"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/docs"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/job"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/middleware"
	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/queue"
	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/metrics"
	redis "github.com/Jeremiah-Williams1/inference-orchestrator/pkg/redisclient"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router holds every dependency the application needs.
// Add new dependencies here as the project grows.
type Router struct {
	cfg   *config.Config
	log   *slog.Logger
	redis *redis.Client
	reg   *prometheus.Registry
}

func New(cfg *config.Config, log *slog.Logger, redisClient *redis.Client, reg *prometheus.Registry) *Router {
	return &Router{cfg: cfg, log: log, redis: redisClient, reg: reg}
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
	m := metrics.NewMetrics(s.reg)

	// Health lives outside /api/v1 — never versioned, always stable.
	r.GET("/health", impl.GetHealth)
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{})))

	spec, err := os.ReadFile("api.yaml")
	if err != nil {
		s.log.Warn("could not read api.yaml — /docs/spec will be unavailable", "err", err)
	}
	docs := docs.NewDocsHandler(spec)
	r.GET("/docs", docs.UI)
	r.GET("/docs/spec", docs.Spec)

	if s.redis != nil {
		redisQueue := queue.NewRedisQueue(s.redis.Redis())
		jobSvc := job.NewJobService(redisQueue, s.log, m)
		impl.job = job.NewJobHandler(jobSvc, s.log)
	}

	v1 := r.Group("/api/v1")
	api.RegisterHandlersWithOptions(v1, impl, api.GinServerOptions{})

	return r
}
