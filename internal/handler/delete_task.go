package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid task id",
		})
	}

	err = h.repo.DeleteTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to delete",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
