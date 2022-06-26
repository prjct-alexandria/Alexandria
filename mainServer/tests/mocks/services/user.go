package services

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

type UserServiceMock struct {
	Mock *mocks.Mock
}

func NewUserServiceMock() UserServiceMock {
	return UserServiceMock{Mock: mocks.NewMock()}
}

var SaveUserMock func(user entities.User) error

func (m UserServiceMock) SaveUser(user entities.User) error {
	m.Mock.CallFunc("SaveUser", &map[string]interface{}{
		"user": user,
	})
	return SaveUserMock(user)
}

var CheckPasswordMock func(email string, pwdClaim string) (entities.User, error)

func (m UserServiceMock) CheckPassword(email string, pwdClaim string) (entities.User, error) {
	m.Mock.CallFunc("CheckPassword", &map[string]interface{}{
		"email":    email,
		"pwdClaim": pwdClaim,
	})
	return CheckPasswordMock(email, pwdClaim)
}
