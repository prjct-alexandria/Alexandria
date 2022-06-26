package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services/interfaces"
	"mainServer/utils/auth"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type RequestController struct {
	Serv interfaces.RequestService
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
	// Check if user is logged in
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

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
	req, err := contr.Serv.CreateRequest(article, form.SourceVersionID, form.TargetVersionID, loggedInAs)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed accepting request on server"))
		return
	}

	c.Header("Content-Type", "application/json")
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
	// Check if user is logged in
	// Checking if this owner is allowed to reject is done by the service
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

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

	// reject request with service
	err = contr.Serv.RejectRequest(request, loggedInAs)
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
	// Check if user is logged in
	// Checking if this owner is allowed to accept is done by the service
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

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
	// accept request with service
	err = contr.Serv.AcceptRequest(request, loggedInAs)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed accepting request on server"))
		return
	}

	c.Status(http.StatusOK)
}

// GetRequest godoc
// @Summary     Get Request
// @Description Returns the information of a given request, including the information of both versions. Note that comparing target and source versions directly, isn't reliable as before-and-after comparison. That's why, instead of filling in the contents of the version fields, a before and after string is included in the response.
// @Accept      plain
// @Produce 	json
// @Param		articleID	path	string	true "Article ID"
// @Param		requestID	path	string	true "Request ID"
// @Success     200 {object} models.RequestWithComparison
// @Failure     400 {object} httperror.HTTPError
// @Failure     500 {object} httperror.HTTPError
// @Router      /articles/{articleID}/requests/{requestID} [get]
func (contr RequestController) GetRequest(c *gin.Context) {
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

	// accept request with service
	req, err := contr.Serv.GetRequest(request)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed getting request on server"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, req)
}

// GetRequestList 	godoc
// @Summary		Get a list of merge requests
// @Description	Gets a list of merge requests (with possible filtering conditions)
// @Param		articleID	path	string	true	"Article ID"
// @Param		sourceID	query	string	false	"Source version"
// @Param		targetID	query	string	false	"Target version"
// @Param		relatedID	query	string	false	"Source or Target version"
// @Produce		json
// @Success		200 {object} []models.Request
// @Failure		400 "Invalid article ID provided"
// @Failure		404 "Could not find merge requests for this article"
// @Router		/articles/{articleID}/requests [get]
func (contr RequestController) GetRequestList(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	aid := c.Param("articleID")
	articleId, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	var sourceId int64
	sid := c.Query("sourceID")
	if sid != "" {
		sourceId, err = strconv.ParseInt(sid, 10, 64)
		if err != nil {
			fmt.Println(err)
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid source ID, cannot interpret as integer, id=%s ", sid))
			return
		}
	} else {
		sourceId = -1
	}

	var targetId int64
	tid := c.Query("targetID")
	if tid != "" {
		targetId, err = strconv.ParseInt(tid, 10, 64)
		if err != nil {
			fmt.Println(err)
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid target ID, cannot interpret as integer, id=%s ", tid))
			return
		}
	} else {
		targetId = -1
	}

	var relatedId int64
	rid := c.Query("relatedID")
	if rid != "" {
		relatedId, err = strconv.ParseInt(rid, 10, 64)
		if err != nil {
			fmt.Println(err)
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid related ID, cannot interpret as integer, id=%s ", rid))
			return
		}
	} else {
		relatedId = -1
	}

	list, err := contr.Serv.GetRequestList(articleId, sourceId, targetId, relatedId)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("could not fetch request list"))
		return
	}

	c.JSON(http.StatusOK, list)
}
