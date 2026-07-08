package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *TaskHandler) GeatAllTasks(c *gin.Context) {
	tasks, err := h.repo.GeatAllTasks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no tasks found",
		})
	}

	c.JSON(http.StatusOK, tasks)
}
