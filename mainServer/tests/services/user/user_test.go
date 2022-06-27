package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/entities"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"testing"
)

var userRepoMock repositories.UserRepositoryMock

var serv *services.UserService

func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test in this file starts
func globalSetup() {

}

// localSetup should be called before each individual test
func localSetup() {
	// Make a clean controller with clean mocks
	userRepoMock = repositories.NewUserRepositoryMock()
	servVal := services.NewUserService(userRepoMock)
	serv = &servVal
}

func TestSaveUserSuccess(t *testing.T) {
	// Arrange
	localSetup()

	user := entities.User{
		Name:  "johnDoe",
		Email: "johnDoe@gmail.com",
		Pwd:   "password",
	}

	repositories.CreateUserMock = func(user entities.User) error {
		return nil
	}

	// Act
	err := serv.SaveUser(user)

	// Assert
	assert.Equal(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "CreateUser", &map[string]interface{}{
		"user": user,
	})
}

func TestSaveUserFail(t *testing.T) {
	// Arrange
	localSetup()

	user := entities.User{
		Name:  "johnDoe",
		Email: "johnDoe@gmail.com",
		Pwd:   "password",
	}

	repositories.CreateUserMock = func(user entities.User) error {
		return errors.New("error")
	}

	// Act
	err := serv.SaveUser(user)

	// Assert
	assert.NotEqual(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "CreateUser", &map[string]interface{}{
		"user": user,
	})
}

func TestCheckPassWordSuccess(t *testing.T) {
	// Arrange
	localSetup()

	emailAdr := "johnDoe@gmail.com"
	pwdClaim := "password"

	user, _ := entities.User{
		Name:  "johnDoe",
		Email: emailAdr,
		Pwd:   "password",
	}.Hash()

	repositories.GetFullUserByEmailMock = func(email string) (entities.User, error) {
		if email == emailAdr {
			return user, nil
		}
		return entities.User{}, nil
	}

	expected := user

	// Act
	actual, err := serv.CheckPassword(emailAdr, pwdClaim)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "GetFullUserByEmail", &map[string]interface{}{
		"email": emailAdr,
	})
}

func TestCheckPasswordDbFail(t *testing.T) {
	// Arrange
	localSetup()

	emailAdr := "johnDoe@gmail.com"
	pwdClaim := "password"

	user, _ := entities.User{
		Name:  "johnDoe",
		Email: emailAdr,
		Pwd:   "password",
	}.Hash()

	repositories.GetFullUserByEmailMock = func(email string) (entities.User, error) {
		if email == emailAdr {
			return user, errors.New("error")
		}
		return entities.User{}, errors.New("error")
	}

	expected := entities.User{}

	// Act
	actual, err := serv.CheckPassword(emailAdr, pwdClaim)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "GetFullUserByEmail", &map[string]interface{}{
		"email": emailAdr,
	})
}

func TestCheckPasswordNonidentical(t *testing.T) {
	// Arrange
	localSetup()

	emailAdr := "johnDoe@gmail.com"
	pwdClaim := "otherPassword"

	user, _ := entities.User{
		Name:  "johnDoe",
		Email: emailAdr,
		Pwd:   "password",
	}.Hash()

	repositories.GetFullUserByEmailMock = func(email string) (entities.User, error) {
		if email == emailAdr {
			return user, nil
		}
		return entities.User{}, nil
	}

	expected := user

	// Act
	actual, err := serv.CheckPassword(emailAdr, pwdClaim)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "GetFullUserByEmail", &map[string]interface{}{
		"email": emailAdr,
	})
}

func TestGetUserByEmailSuccess(t *testing.T) {
	// Arrange
	localSetup()

	emailAdr := "johnDoe@gmail.com"

	user := entities.User{
		Name:  "johnDoe",
		Email: emailAdr,
		Pwd:   "password",
	}

	repositories.GetFullUserByEmailMock = func(email string) (entities.User, error) {
		return user, nil
	}

	expected := user

	// Act
	actual, err := serv.GetUserByEmail(emailAdr)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "GetFullUserByEmail", &map[string]interface{}{
		"email": emailAdr,
	})
}

func TestGetUserByEmailFail(t *testing.T) {
	// Arrange
	localSetup()

	emailAdr := "johnDoe@gmail.com"

	user := entities.User{
		Name:  "johnDoe",
		Email: emailAdr,
		Pwd:   "password",
	}

	repositories.GetFullUserByEmailMock = func(email string) (entities.User, error) {
		return user, errors.New("error")
	}

	expected := user

	// Act
	actual, err := serv.GetUserByEmail(emailAdr)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "GetFullUserByEmail", &map[string]interface{}{
		"email": emailAdr,
	})
}
