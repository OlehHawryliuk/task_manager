package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task id",
		})
	}

	task, err := h.repo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task not found",
		})
	}

	c.JSON(http.StatusFound, task)
}
