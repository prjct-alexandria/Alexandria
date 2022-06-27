package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/server/config"
	"mainServer/tests"
	"mainServer/tests/mocks/services"
	"net/http"
	"testing"
)

var userServMock services.UserServiceMock
var contr *controllers.UserController
var r *gin.Engine
var loggedInUser *string

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test in this file starts
func globalSetup() {
	// Setup test router, to test controller endpoints through http
	r = gin.Default()

	// Mock the authentication middleware, that sets the email of logged-in user
	r.Use(func(c *gin.Context) {
		if loggedInUser != nil {
			c.Set("Email", *loggedInUser)
		}
	})

	// route
	r.POST("/users", func(c *gin.Context) {
		contr.Register(c)
	})
	r.POST("/login", func(c *gin.Context) {
		contr.Login(c)
	})
	r.POST("/logout", func(c *gin.Context) {
		contr.Logout(c)
	})

}

// localSetup should be called before each individual test
func localSetup() {
	// Make a clean controller with clean mocks
	userServMock = services.NewUserServiceMock()
	cfg := config.Config{}
	cfg.Auth.JwtSecret = ""
	contr = &controllers.UserController{UserService: userServMock, Cfg: &cfg}
}

func TestRegisterSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.SaveUserMock = func(user entities.User) error {
		return nil
	}

	// set the logged-in user, should be able to register without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/users", exampleRegisterForm,
		http.StatusOK, nil)

	// check the service mock
	userServMock.Mock.AssertCalled(t, "SaveUser", 1)

	// cannot perform normal assertion due to hashed password, so manually compare just the name and email field
	calledWith := userServMock.Mock.Params["SaveUser"]["user"].(entities.User)
	if calledWith.Name != exampleRegisterForm.Name || calledWith.Email != exampleRegisterForm.Email {
		t.Errorf("Expected user to be %v <%v>, but was %v <%v>",
			exampleRegisterForm.Name, exampleRegisterForm.Email, calledWith.Name, calledWith.Email)
	}
}

func TestRegisterBadForm(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.SaveUserMock = func(user entities.User) error {
		return nil
	}

	// set the logged-in user, should be able to register without being logged in
	loggedInUser = nil

	// Make a bad registration form, missing the password
	badForm := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  exampleRegisterForm.Name,
		Email: exampleRegisterForm.Email,
	}

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/users", badForm,
		http.StatusBadRequest, nil)

	// check the service mock
	userServMock.Mock.AssertCalled(t, "SaveUser", 0)
}

func TestRegisterInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.SaveUserMock = func(user entities.User) error {
		return errors.New("fake error for testing")
	}

	// set the logged-in user, should be able to register without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/users", exampleRegisterForm,
		http.StatusConflict, nil)

	// check the service mock
	calledWith := userServMock.Mock.Params["SaveUser"]["user"].(entities.User)
	if calledWith.Name != exampleRegisterForm.Name || calledWith.Email != exampleRegisterForm.Email {
		t.Errorf("Expected user to be %v <%v>, but was %v <%v>",
			exampleRegisterForm.Name, exampleRegisterForm.Email, calledWith.Name, calledWith.Email)
	}
}

// TODO: login
func TestLoginSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CheckPasswordMock = func(email string, pwd string) (entities.User, error) {
		return exampleUser, nil
	}

	// set the logged-in user, should be able to log in without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/login", exampleLoginForm,
		http.StatusOK, exampleLoginResponse)

	// check the service mock
	userServMock.Mock.AssertCalled(t, "CheckPassword", 1)
	userServMock.Mock.AssertCalledWith(t, "CheckPassword", &map[string]interface{}{
		"email":    exampleLoginForm.Email,
		"pwdClaim": exampleLoginForm.Pwd,
	})
}

func TestLoginSuccessAlreadyLoggedIn(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CheckPasswordMock = func(email string, pwd string) (entities.User, error) {
		return exampleUser, nil
	}

	// set the logged-in user, should be able to log in when already logged in
	email := "zoo@mail.com"
	loggedInUser = &email

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/login", exampleLoginForm,
		http.StatusOK, exampleLoginResponse)

	// check the service mock
	userServMock.Mock.AssertCalled(t, "CheckPassword", 1)
	userServMock.Mock.AssertCalledWith(t, "CheckPassword", &map[string]interface{}{
		"email":    exampleLoginForm.Email,
		"pwdClaim": exampleLoginForm.Pwd,
	})
}

func TestLoginBadForm(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CheckPasswordMock = func(email string, pwd string) (entities.User, error) {
		return exampleUser, nil
	}

	// set the logged-in user, should be able to log in without being logged in
	loggedInUser = nil

	// Make a bad login form, missing the password
	badForm := struct {
		Email string `json:"email"`
	}{
		Email: exampleRegisterForm.Email,
	}

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/login", badForm,
		http.StatusBadRequest, nil)

	// check the service mock
	userServMock.Mock.AssertCalled(t, "CheckPassword", 0)
}

func TestLoginInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CheckPasswordMock = func(email string, pwd string) (entities.User, error) {
		return entities.User{}, errors.New("fake error for testing")
	}

	// set the logged-in user, should be able to log in without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/login", exampleLoginForm,
		http.StatusForbidden, nil)

	// check the service mock
	userServMock.Mock.AssertCalled(t, "CheckPassword", 1)
	userServMock.Mock.AssertCalledWith(t, "CheckPassword", &map[string]interface{}{
		"email":    exampleLoginForm.Email,
		"pwdClaim": exampleLoginForm.Pwd,
	})
}

func TestLogoutSuccess(t *testing.T) {
	localSetup()

	// set the logged-in user, should not matter though
	email := "zoo@mail.com"
	loggedInUser = &email

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/logout", nil,
		http.StatusOK, nil)
}

var exampleUser = entities.User{
	Name:  "Maï Naime is com-plex III",
	Email: "john@mail.com",
	Pwd:   "JohnsSecretPassword1234!",
}

var exampleRegisterForm = models.RegisterForm(exampleUser)

var exampleLoginForm = models.LoginForm{
	Email: exampleUser.Email,
	Pwd:   exampleUser.Pwd,
}

var exampleLoginResponse = models.User{
	Name:  exampleUser.Name,
	Email: exampleUser.Email,
}
