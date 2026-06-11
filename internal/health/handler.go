package health

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/response"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	db    pinger
	cache pinger
}

func NewHandler(db pinger, cache pinger) *Handler {
	return &Handler{db: db, cache: cache}
}

// Handle godoc
// @Summary Health check
// @Description Returns liveness status of all configured dependencies
// @Tags health
// @Produce json
// @Success 200 {object} map[string]any
// @Router /health [get]
func (h *Handler) Handle(c *gin.Context) {
	overallStatus := "ok"

	dbStatus := "unconfigured"
	if h.db != nil {
		if err := h.db.Ping(c.Request.Context()); err != nil {
			dbStatus = "error"
			overallStatus = "degraded"
		} else {
			dbStatus = "ok"
		}
	}

	cacheStatus := "unconfigured"
	if h.cache != nil {
		if err := h.cache.Ping(c.Request.Context()); err != nil {
			cacheStatus = "error"
			overallStatus = "degraded"
		} else {
			cacheStatus = "ok"
		}
	}

	response.OK(c, gin.H{
		"status": overallStatus,
		"db":     dbStatus,
		"cache":  cacheStatus,
	})
}
