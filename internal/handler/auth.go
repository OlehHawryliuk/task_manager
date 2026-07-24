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

type AuthHandler struct {
	userRepo *repository.UserRepository
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: repo}
}

// @Summary Register new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 409 {object} apierror.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apierror.NewInvalidRequest(err.Error()))
		return
	}

	existingUser, _ := h.userRepo.GetUserByEmail(req.Email)
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

	if err := h.userRepo.CreateUser(user); err != nil {
		c.Error(apierror.ErrDatabaseError)
		return
	}

	token, err := service.GenerateToken(user.ID.String())
	if err != nil {
		c.Error(apierror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// @Summary Login user
// @Description User login and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 401 {object} apierror.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apierror.NewInvalidRequest(err.Error()))
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		c.Error(apierror.NewUnauthorized("Invalid email or password"))
		return
	}

	if !service.VerifyPassword(user.Password, req.Password) {
		c.Error(apierror.NewUnauthorized("Invalid email or password"))
		return
	}

	token, err := service.GenerateToken(user.ID.String())
	if err != nil {
		c.Error(apierror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}
