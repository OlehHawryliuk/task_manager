package config

import (
	"log"
	"os"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Couldn`t connect to database")
	}

	db.AutoMigrate(&model.Task{}, &model.User{})

	log.Println("db connected successfully")

	DB = db
	return db
}
