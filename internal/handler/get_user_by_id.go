package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := r.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
