package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/entities"
	"mainServer/middlewares"
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
// @Failure		400 "could not read request body"
// @Failure		400 "invalid user JSON provided"
// @Failure		403 "could not generate password hash"
// @Failure		500 "could not save user to database"
// @Router		/users
func (u *UserController) Register(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("could not read request body"))
		return
	}

	var user entities.User
	err = json.Unmarshal(byteBody, &user)
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
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not save user to database"))
		return
	} else {
		c.Status(http.StatusOK)
	}
}

type credentials struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

// Login		godoc
// @Summary		Endpoint for user logging in
// @Description	Takes in user email and password from a JSON body, verifies if they are correct with the database and returns a JWT token
// @Accept		json
// @Success		200 "Success"
// @Failure		400 "could not read request body"
// @Failure		400 "invalid JSON provided
// @Failure		403 "invalid password"
// @Failure		500 "could not create token"
// @Router		/login
func (u *UserController) Login(c *gin.Context) {
	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("could not read request body"))
		return
	}

	var cred credentials
	err = json.Unmarshal(byteBody, &cred)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("invalid JSON provided"))
		return
	}

	//Check email + pwd combo
	err = u.UserService.CheckPassword(cred.Email, cred.Pwd)
	if err != nil {
		httperror.NewError(c, http.StatusForbidden, errors.New("invalid password"))
		return
	}

	err = middlewares.UpdateJwtCookie(c, cred.Email)

	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not create token"))
		return
	}
	c.Status(http.StatusOK)
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
