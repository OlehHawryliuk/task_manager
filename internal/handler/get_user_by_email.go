package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := r.repo.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, user)
}
