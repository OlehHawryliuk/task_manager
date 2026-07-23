package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/OlehHawryliuk/task_manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	repo         *repository.TaskRepository
	userRepo     *repository.UserRepository
	cacheService *service.CacheService
}

func NewTaskHandler(repo *repository.TaskRepository, userRepo *repository.UserRepository, cacheService *service.CacheService) *TaskHandler {
	return &TaskHandler{
		repo:         repo,
		userRepo:     userRepo,
		cacheService: cacheService,
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

	ctx := c.Request.Context()
	h.cacheService.Delete(ctx, "tasks:all")
	c.JSON(http.StatusCreated, task)
}

// @Summary Get a task by ID
// @Description Retrieve a specific task by its unique UUID
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Task UUID"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string "Invalid task ID or task not found"
// @Router /tasks/{id} [get]
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
	ctx := c.Request.Context()
	cacheKey := "tasks:all"

	cachedData, err := h.cacheService.Get(ctx, cacheKey)
	if err == nil {
		c.JSON(http.StatusOK, json.RawMessage(cachedData))
		return
	}

	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no tasks found",
		})

		return
	}

	h.cacheService.Set(ctx, cacheKey, tasks, service.TaskCacheTTl)
	c.JSON(http.StatusOK, tasks)
}

// @Summary Update a task
// @Description Update task data (accessible by owner or admin only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Task UUID"
// @Param request body UpdateTaskRequest true "Task data to update"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: you can only update your own tasks"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
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

	ctx := c.Request.Context()
	h.cacheService.Delete(ctx, "tasks:all")
	c.JSON(http.StatusOK, task)
}

// @Summary Delete a task
// @Description Delete a specific task (accessible by owner or admin only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Task UUID"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 400 {object} map[string]string "Invalid task ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: you can only delete your own tasks"
// @Failure 404 {object} map[string]string "Task not found"
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

	ctx := c.Request.Context()
	h.cacheService.Delete(ctx, "tasks:all")
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
