package handler

import (
	"net/http"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/apierror"
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

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// @Summary Create a new user
// @Description Register a new user in the system with default 'user' role (admin only)
// @Tags Users (Admin)
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param request body CreateUserRequest true "User registration data"
// @Success 201 {object} model.User
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 409 {object} apierror.ErrorResponse
// @Failure 500 {object} apierror.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apierror.NewInvalidRequest(err.Error()))
		return
	}

	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.Error(apierror.NewConflict("Email already exists"))
		return
	}

	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.Error(apierror.ErrInternalServer)
		return
	}

	user := &model.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.CreateUser(user); err != nil {
		c.Error(apierror.ErrDatabaseError)
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
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 401 {object} apierror.ErrorResponse
// @Failure 403 {object} apierror.ErrorResponse
// @Failure 404 {object} apierror.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	currentUserID := c.GetString("userID")

	parsedCurrentUserID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.Error(apierror.NewInvalidRequest("Invalid user ID"))
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.Error(apierror.NewInvalidRequest("Invalid user ID"))
		return
	}

	currentUser, err := h.repo.GetUserByID(parsedCurrentUserID)
	if err != nil || currentUser == nil {
		c.Error(apierror.ErrUnauthorized)
		return
	}

	if parsedCurrentUserID != userID && currentUser.Role != "admin" {
		c.Error(apierror.NewForbidden("You can only view your own profile"))
		return
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil || user == nil {
		c.Error(apierror.NewNotFound("User not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get all users
// @Description Retrieve a list of all registered users (admin only)
// @Tags Users (Admin)
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.User
// @Failure 401 {object} apierror.ErrorResponse
// @Failure 403 {object} apierror.ErrorResponse
// @Failure 500 {object} apierror.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		c.Error(apierror.ErrDatabaseError)
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get user by email
// @Description Find a specific user by their email address
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param email path string true "User email"
// @Success 200 {object} model.User
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 401 {object} apierror.ErrorResponse
// @Failure 404 {object} apierror.ErrorResponse
// @Router /users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.Error(apierror.NewInvalidRequest("Email is required"))
		return
	}

	user, err := h.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		c.Error(apierror.NewNotFound("User not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user profile
// @Description Update user data (accessible by profile owner or admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User UUID"
// @Param request body UpdateUserRequest true "Updated user data"
// @Success 200 {object} model.User
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 401 {object} apierror.ErrorResponse
// @Failure 403 {object} apierror.ErrorResponse
// @Failure 404 {object} apierror.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	currentID := c.GetString("userID")

	parsedCurrentUserID, err := uuid.Parse(currentID)
	if err != nil {
		c.Error(apierror.NewInvalidRequest("Invalid user ID"))
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.Error(apierror.NewInvalidRequest("Invalid user ID"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apierror.NewInvalidRequest(err.Error()))
		return
	}

	currentUser, err := h.repo.GetUserByID(parsedCurrentUserID)
	if err != nil || currentUser == nil {
		c.Error(apierror.ErrUnauthorized)
		return
	}

	if parsedCurrentUserID != userID && currentUser.Role != "admin" {
		c.Error(apierror.NewForbidden("You can only update your own profile"))
		return
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil || user == nil {
		c.Error(apierror.NewNotFound("User not found"))
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
			c.Error(apierror.ErrInternalServer)
			return
		}
		user.Password = hashedPassword
	}

	if req.Role != "" && currentUser.Role == "admin" {
		user.Role = req.Role
	}

	user.UpdatedAt = time.Now()

	if err := h.repo.UpdateUser(user); err != nil {
		c.Error(apierror.ErrDatabaseError)
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
// @Success 200 {object} map[string]string
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 401 {object} apierror.ErrorResponse
// @Failure 403 {object} apierror.ErrorResponse
// @Failure 404 {object} apierror.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	currentID := c.GetString("userID")

	parsedCurrentID, err := uuid.Parse(currentID)
	if err != nil {
		c.Error(apierror.NewInvalidRequest("Invalid user ID"))
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.Error(apierror.NewInvalidRequest("Invalid user ID"))
		return
	}

	currentUser, err := h.repo.GetUserByID(parsedCurrentID)
	if err != nil || currentUser == nil {
		c.Error(apierror.ErrUnauthorized)
		return
	}

	if userID != parsedCurrentID && currentUser.Role != "admin" {
		c.Error(apierror.NewForbidden("You can only delete your own profile"))
		return
	}

	targetUser, err := h.repo.GetUserByID(userID)
	if err != nil || targetUser == nil {
		c.Error(apierror.NewNotFound("User not found"))
		return
	}

	if err := h.repo.DeleteUser(userID); err != nil {
		c.Error(apierror.ErrDatabaseError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
