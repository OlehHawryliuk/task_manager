package repository

import (
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTask(c *gin.Context) {
	task := model.Task{
		ID:          uuid.New(),
		Title:       "make new project",
		Description: "use gorm and gin",
		Done:        false,
		UserID:      uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := config.DB.Create(&task)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"task": task,
	})
}
