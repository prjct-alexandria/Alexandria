package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/middlewares"
	"mainServer/models"
	"mainServer/server/config"
	"mainServer/services/interfaces"
	"mainServer/utils/httperror"
	"net/http"
)

type UserController struct {
	UserService interfaces.UserService
	Cfg         *config.Config
}

// Register		godoc
// @Summary		Endpoint for user registration
// @Description	Takes in user credentials from a JSON body, and makes sure they are securely stored in the database.
// @Accept		json
// @Param		credentials	body	entities.User true "User credentials"
// @Success		200 "Success"
// @Failure		400 "Could not read request body"
// @Failure		400 "Invalid user JSON provided"
// @Failure		403 "Could not generate password hash"
// @Failure		409 "Could not save user to database"
// @Router		/users	[post]
func (u *UserController) Register(c *gin.Context) {
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
// @Param		credentials	body	models.LoginForm true "User credentials"
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

	err = middlewares.UpdateJwtCookie(c, cred.Email, u.Cfg)

	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not create token"))
		return
	}
	c.IndentedJSON(http.StatusOK, models.User{Name: dbUser.Name, Email: dbUser.Email})
}

// Logout		godoc
// @Summary		Endpoint for user logging out
// @Description	Sets an expired cookie with an empty email and returns a JWT token
// @Accept		json
// @Param		credentials	body	models.LoginForm true "User credentials"
// @Success		200 "Success"
// @Failure		400 "Could not read request body"
// @Failure		400 "Invalid JSON provided"
// @Failure		500 "Could not update token"
// @Router		/logout	[post]
func (u *UserController) Logout(c *gin.Context) {
	err := middlewares.ExpireJwtCookie(c, u.Cfg)

	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not delete token to logout user"))
		return
	}

	c.IndentedJSON(http.StatusOK, models.User{Name: "", Email: ""})
}
