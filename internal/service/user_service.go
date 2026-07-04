package service

import "github.com/medhansh-32/api-gateway/internal/models"

type UserService interface {
	getUserById(userId int64) (models.User)
}

