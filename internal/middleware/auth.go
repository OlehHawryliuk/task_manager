package middleware

import (
	"net/http"
	"strings"

	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/OlehHawryliuk/task_manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	repo *repository.UserRepository
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
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

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString("userID")
		if userID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Misiing User ID",
			})
			ctx.Abort()
			return
		}

		ParsedUserID, err := uuid.Parse(userID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID",
			})
			ctx.Abort()
			return
		}

		user, err := config.UserRepo.GetUserByID(ParsedUserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			ctx.Abort()
			return
		}

		if user.Role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden: Admin access required",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
