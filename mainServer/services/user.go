package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mainServer/entities"
	"mainServer/repositories/interfaces"
)

type UserService struct {
	UserRepository interfaces.UserRepository
}

func (u *UserService) SaveUser(user entities.User) error {
	return u.UserRepository.CreateUser(user)
}

func (u *UserService) CheckPassword(email string, pwdClaim string) error {
	dbUser, err := u.UserRepository.GetFullUserByEmail(email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(dbUser.Pwd), []byte(pwdClaim))
}

func (u *UserService) GetUserByEmail(email string) (entities.User, error) {
	return u.UserRepository.GetFullUserByEmail(email)
}
