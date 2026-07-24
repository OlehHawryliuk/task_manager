package middleware

import (
	"net/http"

	"github.com/OlehHawryliuk/task_manager/internal/apierror"
	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			if apiErr, ok := err.Err.(*apierror.APIError); ok {
				c.JSON(apiErr.StatusCode, apierror.ErrorResponse{
					Code:    apiErr.Code,
					Message: apiErr.Message,
					Details: apiErr.Details,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, apierror.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal server error",
				Details: err.Error(),
			})
		}
	}
}
