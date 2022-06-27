package repositories

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

type UserRepositoryMock struct {
	Mock *mocks.Mock
}

func NewUserRepositoryMock() UserRepositoryMock {
	return UserRepositoryMock{Mock: mocks.NewMock()}
}

var CreateUserMock func(user entities.User) error

func (m UserRepositoryMock) CreateUser(user entities.User) error {
	m.Mock.CallFunc("CreateUser", &map[string]interface{}{
		"user": user,
	})
	return CreateUserMock(user)
}

var GetFullUserByEmailMock func(email string) (entities.User, error)

func (m UserRepositoryMock) GetFullUserByEmail(email string) (entities.User, error) {
	m.Mock.CallFunc("GetFullUserByEmail", &map[string]interface{}{
		"email": email,
	})
	return GetFullUserByEmailMock(email)
}

var UpdateUserMock func(email string, user entities.User) error

func (m UserRepositoryMock) UpdateUser(email string, user entities.User) error {
	m.Mock.CallFunc("UpdateUser", &map[string]interface{}{
		"email": email,
		"user":  user,
	})
	return UpdateUserMock(email, user)
}

var DeleteUserMock func(email string) error

func (m UserRepositoryMock) DeleteUser(email string) error {
	m.Mock.CallFunc("DeleteUser", &map[string]interface{}{
		"email": email,
	})
	return DeleteUserMock(email)
}

var CheckIfExistsMock func(email string) (bool, error)

func (m UserRepositoryMock) CheckIfExists(email string) (bool, error) {
	m.Mock.CallFunc("CheckIfExists", &map[string]interface{}{
		"email": email,
	})
	return CheckIfExistsMock(email)
}
