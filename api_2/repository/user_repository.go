package repository

import (
	"errors"
	"main/model"

	"gorm.io/gorm"
)

var _ IUserRepository = &UserRepository{}

type IUserRepository interface {
	GetUserByID(id uint) (*model.User, bool, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByID(id uint) (*model.User, bool, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, errors.Is(err, gorm.ErrRecordNotFound), err
	}

	return &user, false, nil
}
