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

// @Summary Create/Register a new user
// @Description Register a new user in the system with default 'user' role
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.User true "User registration data (Email, Username, Password)"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string "Invalid input data or user already exists"
// @Failure 500 {object} map[string]string "Failed to hash password"
// @Router /users [post]
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

// @Summary Get user by ID
// @Description Retrieve profile data for a specific user (accessible by profile owner or admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "Invalid user ID or user not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: you can only view your own profile"
// @Router /users/{id} [get]
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

// @Summary Get all users
// @Description Retrieve a list of all registered users (accessible by admin only via middleware)
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.User
// @Failure 400 {object} map[string]string "Failed to fetch users"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: Admin only"
// @Router /users [get]
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

// @Summary Get user by email
// @Description Find a specific user by their email address (accessible by admin only via middleware)
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param email path string true "User email"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "User not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: Admin only"
// @Router /users/email/{email} [get]
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

// @Summary Update user profile
// @Description Update user data (accessible by profile owner or admin only. Role changes allowed for admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Param request body model.User true "Updated user data"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "Invalid input or user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: You can only update your own profile"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [put]
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
	if req.Role != "" && currentUser.Role == "admin" {
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

// @Summary Delete user account
// @Description Permanently delete a user account (accessible by profile owner or admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID or failed to delete"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden: You can only delete your own profile"
// @Failure 500 {object} map[string]string "User not found"
// @Router /users/{id} [delete]
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
