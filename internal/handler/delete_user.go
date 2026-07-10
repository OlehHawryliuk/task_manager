package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err = r.repo.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete user",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted succesfully",
	})
}
