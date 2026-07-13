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

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
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

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no tasks found",
		})
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	var req UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.repo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	task.Done = req.Done
	task.UpdatedAt = time.Now()

	if err := h.repo.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

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
