package handler

import (
	"net/http"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to parse user ID",
		})
	}

	var req model.User
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	user, err := r.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := r.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
