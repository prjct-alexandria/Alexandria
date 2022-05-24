package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/services"
	"net/http"
)

type UserController struct {
	UserService services.UserService
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
		c.String(http.StatusOK, "Success")
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
