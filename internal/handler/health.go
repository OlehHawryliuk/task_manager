package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
	Version  string `json:"version"`
	Time     string `json:"time"`
}

// @Summary Health Check
// @Description Check API, database, and Redis status
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	health := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := config.DB.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		health.Database = "unhealthy"
		health.Status = "degraded"
	} else {
		health.Database = "healthy"
	}

	if config.RedisClient != nil {
		_, err := config.RedisClient.Ping(ctx).Result()
		if err == nil {
			health.Redis = "healthy"
		} else {
			health.Redis = "unhealthy"
			health.Status = "degraded"
		}
	} else {
		health.Redis = "not_configured"
	}

	statusCode := http.StatusOK
	if health.Status != "degraded" {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}
