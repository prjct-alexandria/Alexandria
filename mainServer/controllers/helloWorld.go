package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelloWorldController struct{}

func (contr HelloWorldController) GetHelloWorld(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}
