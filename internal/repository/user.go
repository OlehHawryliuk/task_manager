package repository

import (
	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoty struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoty {
	return &UserRepositoty{db: db}
}

func (r *UserRepositoty) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoty) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoty) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "email = ?", email).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoty) GetAllUsers() ([]model.User, error) {
	var users []model.User

	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoty) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoty) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}
