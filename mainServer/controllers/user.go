package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/middlewares"
	"mainServer/models"
	"mainServer/services"
	"mainServer/utils/httperror"
	"net/http"
)

type UserController struct {
	UserService services.UserService
}

// Register		godoc
// @Summary		Endpoint for user registration
// @Description	Takes in user credentials from a JSON body, and makes sure they are securely stored in the database.
// @Accept		json
// @Success		200 "Success"
// @Failure		400 "Could not read request body"
// @Failure		400 "Invalid user JSON provided"
// @Failure		403 "Could not generate password hash"
// @Failure		409 "Could not save user to database"
// @Router		/users	[post]
func (u *UserController) Register(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	var user entities.User

	err := c.BindJSON(&user)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("invalid user JSON provided"))
		return
	}

	hashedUser, err := user.Hash()
	if err != nil {
		httperror.NewError(c, http.StatusForbidden, errors.New("could not generate password hash"))
		return
	}
	err = u.UserService.SaveUser(hashedUser)

	if err != nil {
		httperror.NewError(c, http.StatusConflict, errors.New("email address already in use"))
		return
	} else {
		c.Status(http.StatusOK)
	}
}

// Login		godoc
// @Summary		Endpoint for user logging in
// @Description	Takes in user email and password from a JSON body, verifies if they are correct with the database and returns a JWT token
// @Accept		json
// @Success		200 "Success"
// @Failure		400 "Could not read request body"
// @Failure		400 "Invalid JSON provided"
// @Failure		403 "Invalid password"
// @Failure		500 "Could not create token"
// @Router		/login	[post]
func (u *UserController) Login(c *gin.Context) {
	var cred models.LoginForm

	err := c.BindJSON(&cred)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("invalid JSON provided"))
		return
	}

	//Check email + pwd combo
	dbUser, err := u.UserService.CheckPassword(cred.Email, cred.Pwd)
	if err != nil {
		httperror.NewError(c, http.StatusForbidden, errors.New("invalid email and password combination"))
		return
	}

	err = middlewares.UpdateJwtCookie(c, cred.Email)

	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not create token"))
		return
	}
	c.IndentedJSON(http.StatusOK, models.User{Name: dbUser.Name, Email: dbUser.Email})
}

// CreateExampleUser godoc
// @Summary      Temporary user creation endpoint
// @Description  Creates a hardcoded user entity and adds it to the database, demonstrates how to add to database
// @Produce      plain
// @Success      200  {object} string
// @Router       /createExampleUser [post]
func (u *UserController) CreateExampleUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	user := entities.User{
		Name:  "Pietje",
		Email: "pietjepuk@gmail.com",
		Pwd:   "password123",
	}

	err := u.UserService.SaveUser(user)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "Fail")
	} else {
		c.String(http.StatusOK, "Succes")
	}
}

// GetExampleUser godoc
// @Summary      Get test user from database endpoint
// @Description  Returns a user with a hardcoded email address, demonstrates how to use the services.
// @Produce      json
// @Success      200  {object} entities.User
// @Router       /getExampleUser [get]
func (u *UserController) GetExampleUser(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	user, err := u.UserService.GetUserByEmail("pietjepuk@gmail.com")

	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "Fail")
	} else {
		c.IndentedJSON(http.StatusOK, user)
	}
}
