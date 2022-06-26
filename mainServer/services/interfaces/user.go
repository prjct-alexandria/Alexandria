package interfaces

import "mainServer/entities"

type UserService interface {
	SaveUser(user entities.User) error
	CheckPassword(email string, pwdClaim string) (entities.User, error)
}
