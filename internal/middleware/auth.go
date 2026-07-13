package middleware

import (
	"net/http"
	"strings"

	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/OlehHawryliuk/task_manager/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repo *repository.UserRepository
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Autorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Token not found",
			})

			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Invalid token format",
			})

			ctx.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := service.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})

			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Next()
	}
}
