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

// CreateRequest godoc
// @Summary     Create request
// @Description Creates request to merge one article versions' changes into another
// @Accept      json
// @Produce 	json
// @Param		articleID	path	string	true "Article ID"
// @Param		request		body	models.RequestCreationForm	true "Request"
// @Success     200 {object} models.Request
// @Failure     400 "Invalid article ID or request creation data"
// @Failure     500 "Error creating request on server"
// @Router      /articles/{articleID}/requests [post]
func (contr RequestController) CreateRequest(c *gin.Context) {
	c.Header("Content-Type", "application/json")

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
	req, err := contr.Serv.CreateRequest(article, form.SourceVersionID, form.TargetVersionID)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed creating article on server"))
		return
	}

	c.JSON(http.StatusOK, req)
}

// RejectRequest godoc
// @Summary     Reject request
// @Description Rejects request to merge one article versions' changes into another.
// @Accept      plain
// @Produce 	json
// @Param		articleID	path	string	true "Article ID"
// @Param		requestID	path	string	true "Request ID"
// @Success     200
// @Failure     400 {object} httperror.HTTPError
// @Failure     403 {object} httperror.HTTPError
// @Failure     500 {object} httperror.HTTPError
// @Router      /articles/{articleID}/requests/{requestID}/reject [put]
func (contr RequestController) RejectRequest(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// extract article id, had it in the path for consistency in endpoints, but actually ignores it
	aid := c.Param("articleID")
	_, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// extract request id
	rid := c.Param("requestID")
	request, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid version ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// get the current logged-in user
	val, exists := c.Get("Email")
	email := fmt.Sprintf("%v", val) // convert email interface{} type to string
	if !exists {
		fmt.Println(err)
		httperror.NewError(c, http.StatusUnauthorized, fmt.Errorf("you have to be logged-in to reject a request"))
		return
	}

	// reject request with service
	err = contr.Serv.RejectRequest(request, email)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed rejecting request on server"))
		return
	}

	c.Status(http.StatusOK)
}

// AcceptRequest godoc
// @Summary     Accepts request
// @Description Accepts request to merge one article versions' changes into another. Updates target version and archives the request, by recording the current latest commits and setting its state to 'accepted'.
// @Accept      plain
// @Produce 	json
// @Param		articleID	path	string	true "Article ID"
// @Param		requestID	path	string	true "Request ID"
// @Success     200
// @Failure     400 {object} httperror.HTTPError
// @Failure     403 {object} httperror.HTTPError
// @Failure     500 {object} httperror.HTTPError
// @Router      /articles/{articleID}/requests/{requestID}/accept [put]
func (contr RequestController) AcceptRequest(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// extract article id, had it in the path for consistency in endpoints, but actually ignores it
	aid := c.Param("articleID")
	_, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// extract request id
	rid := c.Param("requestID")
	request, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid version ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// get the current logged-in user
	val, exists := c.Get("Email")
	if val == nil || exists == false {
		fmt.Println(err)
		httperror.NewError(c, http.StatusUnauthorized, fmt.Errorf("you have to be logged-in to accepting a request"))
		return
	}
	email := fmt.Sprintf("%v", val) // convert email interface{} type to string

	// accept request with service
	err = contr.Serv.AcceptRequest(request, email)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed accepting request on server"))
		return
	}

	c.Status(http.StatusOK)
}
