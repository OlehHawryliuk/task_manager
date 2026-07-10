package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
