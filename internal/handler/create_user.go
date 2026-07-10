package handler

import (
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	repo *repository.UserRepositoty
}

func NewUserHandler(repo *repository.UserRepositoty) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h UserHandler) CreateUser(c *gin.Context) {
	var req model.User
	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := &model.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	err = h.repo.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}
