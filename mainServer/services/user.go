package services

import (
	"mainServer/entities"
	"mainServer/repositories/interfaces"
)

type UserService struct {
	UserRepository interfaces.UserRepository
}

func (u *UserService) SaveUser(user entities.User) error {
	return u.UserRepository.CreateUser(user)
}

func (u *UserService) GetUserByEmail(email string) (entities.User, error) {
	return u.UserRepository.GetFullUserByEmail(email)
}
