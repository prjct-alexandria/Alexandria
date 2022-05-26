package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type RequestController struct {
	Serv services.RequestService
}

func (contr RequestController) CreateRequest(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	// read path parameter
	aid := c.Param("articleID")
	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// read request creation JSON
	form := models.RequestCreationForm{}
	err = c.BindJSON(&form)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, errors.New("cannot bind request creation form"))
		return
	}

	// create request with service
	req, err := contr.Serv.CreateRequest(article, form.SourceVersionID, form.TargetVersionID, form.SourceHistoryID, form.TargetHistoryID)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed creating article on server"))
		return
	}

	c.JSON(http.StatusOK, req)
}
