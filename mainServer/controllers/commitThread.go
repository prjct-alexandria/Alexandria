package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services"
	"net/http"
)

type CommitThreadController struct {
	CommitThreadService services.CommitThreadService
}

// TODO: make 2 different endpoints for creating a thread and creating a specific one
func (contr *CommitThreadController) CreateThread(c *gin.Context) {
	thread := models.CommitThreadNoId{}
	err := c.BindJSON(&thread)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	aid := c.Param("articleID")
	cid := c.Param("commitID")

	if err := contr.CommitThreadService.StartThread(c, thread, aid, cid); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Status(http.StatusOK)
}
