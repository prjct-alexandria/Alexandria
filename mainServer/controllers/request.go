package controllers

import (
	"github.com/gin-gonic/gin"
	"mainServer/services"
)

type RequestController struct {
	Serv services.RequestService
}

func (contr RequestController) CreateRequest(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	aid := c.Param("articleID")
	vid := c.Param("versionID")

}
