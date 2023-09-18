package service

import (
	"main/model"
	"main/repository"
)

var _ IUserService = &UserService{}

type IUserService interface {
	GetUserByID(id uint) (*model.User, bool, error)
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) GetUserByID(id uint) (*model.User, bool, error) {
	return us.userRepository.GetUserByID(id)
}
