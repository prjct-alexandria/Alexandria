package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mainServer/entities"
	"mainServer/repositories/interfaces"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(userRepository interfaces.UserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func (u UserService) SaveUser(user entities.User) error {
	return u.userRepository.CreateUser(user)
}

func (u UserService) CheckPassword(email string, pwdClaim string) (entities.User, error) {
	dbUser, err := u.userRepository.GetFullUserByEmail(email)
	if err != nil {
		fmt.Println(err)
		return entities.User{}, err
	}
	return dbUser, bcrypt.CompareHashAndPassword([]byte(dbUser.Pwd), []byte(pwdClaim))
}

func (u UserService) GetUserByEmail(email string) (entities.User, error) {
	return u.userRepository.GetFullUserByEmail(email)
}
