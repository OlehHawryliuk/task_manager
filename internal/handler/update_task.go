package handler

import (
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task id",
		})
	}

	var req UpdateTaskRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	task := model.Task{
		ID:          taskID,
		Description: req.Description,
		Title:       req.Title,
		Done:        req.Done,
		UpdatedAt:   time.Now(),
	}

	err = h.repo.UpdateTask(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update task",
		})
	}
	c.JSON(http.StatusOK, task)
}
