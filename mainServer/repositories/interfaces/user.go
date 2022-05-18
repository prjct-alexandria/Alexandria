package interfaces

import (
	"mainServer/entities"
)

type UserRepository interface {
	CreateUser(user entities.User) error
	GetFullUserByEmail(email string) (entities.User, error)
	UpdateUser(email string, user entities.User) error
	DeleteUser(email string) error
}
