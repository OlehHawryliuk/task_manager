package handler

import (
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create task",
		})
	}

	task := model.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Done:        false,
		UserID:      uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.repo.CreateTask(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create task",
		})
	}

	c.JSON(http.StatusCreated, task)
}
