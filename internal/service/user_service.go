package service

import (
	"github.com/medhansh-32/api-gateway/internal/models"
	"github.com/medhansh-32/api-gateway/internal/repository"
)

type UserService interface {
	GetUserById(userId int64) (*models.User, error)
	GetUserByUserName(userName string) (*models.User, error)
}

type UserServiceImpl struct {
	UserRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) (UserService){
	return &UserServiceImpl{UserRepository: userRepository}
}

func (userServiceImpl *UserServiceImpl) GetUserById(userId int64) (*models.User, error) {
	user, err := userServiceImpl.UserRepository.GetUserByUserId(userId)
	if err != nil {
		return nil, err
	}
	return user,nil
}

func (userServiceImpl *UserServiceImpl) GetUserByUserName(userName string) (*models.User, error) {
	user, err := userServiceImpl.UserRepository.GetUserByUserName(userName)
	if err != nil {
		return nil, err
	}
	return user,nil
}
