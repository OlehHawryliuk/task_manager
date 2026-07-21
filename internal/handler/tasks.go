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
	repo     *repository.TaskRepository
	userRepo *repository.UserRepository
}

func NewTaskHandler(repo *repository.TaskRepository, userRepo *repository.UserRepository) *TaskHandler {
	return &TaskHandler{
		repo:     repo,
		userRepo: userRepo,
	}
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

// @Summary Create a new task
// @Description Create a new task for authenticated user
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param request body CreateTaskRequest true "Task data"
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing user ID",
		})

		return
	}

	ParsedUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})

		return
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create task",
		})

		return
	}

	task := model.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Done:        false,
		UserID:      ParsedUserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.repo.CreateTask(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create task",
		})

		return
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

		return
	}

	task, err := h.repo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task not found",
		})

		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary Get all tasks
// @Description Retrieve all tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.Task
// @Failure 401 {object} map[string]string
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no tasks found",
		})

		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary Update a task
// @Description Update task (owner or admin only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Task ID"
// @Param request body UpdateTaskRequest true "Task data"
// @Success 200 {object} model.Task
// @Failure 403 {object} map[string]string
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing user ID",
		})

		return
	}

	ParsedUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})

		return
	}

	user, err := h.userRepo.GetUserByID(ParsedUserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})

		return
	}

	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task id",
		})
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

	if task.UserID != ParsedUserID && user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbiden: you can only update your own tasks",
		})
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

// @Summary Delete a task
// @Description Delete task (owner or admin only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing user ID",
		})

		return
	}

	ParsedUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})

		return
	}

	user, err := h.userRepo.GetUserByID(ParsedUserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})

		return
	}

	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task id",
		})

		return
	}

	task, err := h.repo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})

		return
	}

	if task.UserID != ParsedUserID && user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbiden: you can only delete your own tasks",
		})
		return
	}

	err = h.repo.DeleteTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to delete",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
