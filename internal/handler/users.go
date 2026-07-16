package handler

import (
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/OlehHawryliuk/task_manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.User
	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	password, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := &model.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  password,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (r *UserHandler) GetUserByID(c *gin.Context) {
	currentUserID := c.GetString("userID")
	ParsedCurrentUserID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid userID",
		})
		return
	}

	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	currentUser, err := r.repo.GetUserByID(ParsedCurrentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	if ParsedCurrentUserID != userID && currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: you can only view ypur own profile",
		})
		return
	}

	user, err := r.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := r.repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch users",
		})

		return
	}

	c.JSON(http.StatusOK, users)
}

func (r *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := r.repo.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *UserHandler) UpdateUser(c *gin.Context) {
	currentID := c.GetString("userID")
	if currentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing user ID",
		})
		return
	}

	parsedCurrentUserID, err := uuid.Parse(currentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to parse user ID",
		})
		return
	}

	var req model.User
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	currentUser, err := r.repo.GetUserByID(parsedCurrentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if parsedCurrentUserID != userID && currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You can only update your own profile",
		})
		return
	}

	user, err := r.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := service.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		user.Password = hashedPassword
	}
	if req.Role != "" && user.Role == "admin" {
		user.Role = req.Role
	}

	if err := r.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *UserHandler) DeleteUser(c *gin.Context) {
	currentID := c.GetString("userID")
	ParsedCurrentID, err := uuid.Parse(currentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := r.repo.GetUserByID(ParsedCurrentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	if userID != ParsedCurrentID && user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You can only delete your own profile",
		})
		return
	}

	err = r.repo.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted succesfully",
	})
}
