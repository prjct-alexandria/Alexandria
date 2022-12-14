package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/services/interfaces"
	"mainServer/utils/auth"
	gitUtils "mainServer/utils/git"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type ThreadController struct {
	ThreadService                interfaces.ThreadService
	CommitThreadService          interfaces.CommitThreadService
	CommitSelectionThreadService interfaces.CommitSelectionThreadService
	RequestThreadService         interfaces.RequestThreadService
	ReviewThreadService          interfaces.ReviewThreadService
	CommentService               interfaces.CommentService
}

// CreateThread godoc
// @Summary      Creates thread entity
// @Description  Creates thread entity, and specific thread entity. Returns id's of thread, specific thread and comment
// @Accept		 json
// @Param		 thread 		body	models.Thread		true 	"Thread"
// @Produce      json
// @Success      200  {object} models.ReturnThreadIds
// @Failure 	 400  {object} httperror.HTTPError
// @Failure 	 500  {object} httperror.HTTPError
// @Router       /articles/:articleID/thread/:threadType/id/:specificID/ [post]
func (contr *ThreadController) CreateThread(c *gin.Context) {
	// Check if user is logged in
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

	var thread models.Thread
	err := c.BindJSON(&thread)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	aid := c.Param("articleID")
	sid := c.Param("specificID")
	threadType := c.Param("threadType")

	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// save thread in the db
	tid, err := contr.ThreadService.StartThread(thread, intAid)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// save first comment in the db
	coid, err := contr.CommentService.SaveComment(thread.Comments[0], tid, loggedInAs)

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// If needed, these could be split up over 4 different endpoints instead of using one with multiple responsibilities
	var id int64
	var threadError error

	switch threadType {
	case "commit":
		// check if the specific thread ID string can actually be a commit ID
		if !gitUtils.IsCommitHash(sid) {
			err := fmt.Errorf("invalid commit id=%s, should be a 40-character long hex string", sid)
			httperror.NewError(c, http.StatusBadRequest, err)
			return
		}
		id, threadError = contr.CommitThreadService.StartCommitThread(sid, tid)
	case "commitSelection":
		// check if the specific thread ID string can actually be a commit ID
		if !gitUtils.IsCommitHash(sid) {
			err := fmt.Errorf("invalid commit id=%s, should be a 40-character long hex string", sid)
			httperror.NewError(c, http.StatusBadRequest, err)
			return
		}

		selection := thread.Selection
		if len(selection) > 255 || len(selection) < 1 {
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("invalid selection length, got %d", len(selection)))
			return
		}

		id, threadError = contr.CommitSelectionThreadService.StartCommitSelectionThread(sid, tid, selection)
	case "request":
		intSid, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("invalid requestID, got %v", sid))
			return
		}
		id, threadError = contr.RequestThreadService.StartRequestThread(intSid, tid, loggedInAs)
	case "review":
		intSid, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("invalid reviewID, got %v", sid))
			return
		}
		id, threadError = contr.ReviewThreadService.StartReviewThread(intSid, tid)
	default:
		id, threadError = -1, errors.New("invalid thread type")
	}

	if threadError != nil {
		//TODO: Distinguish between different error types
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not create thread"))
		fmt.Println(err)
		return
	}

	// return tid and specific id
	ids := models.ReturnThreadIds{
		ThreadId:  tid,
		CommentId: coid,
		Id:        id,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, ids)
}

// SaveComment godoc
// @Summary     Save comment
// @Description Save all types (commit/request/review) of comments to the database
// @Accept      json
// @Param 		comment 		body	entities.Comment		true 	"Comment"
// @Param		threadID		path	string					true	"Thread ID"
// @Success     200 "Success"
// @Failure     400 "Bad request"
// @Failure     500 "failed saving comment"
// @Router      /comments/thread/:threadID [post]
func (contr *ThreadController) SaveComment(c *gin.Context) {
	// Check if user is logged in
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

	var comment entities.Comment
	err := c.BindJSON(&comment)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	tid := c.Param("threadID")
	intTid, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("could not parse threadID %v", tid))
		return
	}

	id, err := contr.CommentService.SaveComment(comment, intTid, loggedInAs)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed saving comment"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.IndentedJSON(http.StatusOK, id)
}

// GetRequestThreads godoc
// @Summary      Get all threads for a request
// @Description  Gets a list with all threads belonging to a specific request of an article
// @Param		 article ID 		path	int64		true 	"Article ID"
// @Param		 request ID 		path	int64		true 	"Request ID"
// @Produce      json
// @Success      200  {object} []models.Thread
// @Failure      400  {object} httperror.HTTPError
// @Router       /articles/:articleID/requests/:requestID/threads [get]
func (contr *ThreadController) GetRequestThreads(c *gin.Context) {
	aid := c.Param("articleID")
	rid := c.Param("requestID")

	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	intRid, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	threads, err := contr.RequestThreadService.GetRequestThreads(intAid, intRid)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, fmt.Errorf("cannot find requestthreads for this request"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, threads)
}

// GetCommitThreads godoc
// @Summary      Get all threads for a commit
// @Description  Gets a list with all threads belonging to a specific commit of an article
// @Param		 article ID 		path	int64		true 	"Article ID"
// @Param		 commit ID 			path	int64		true 	"Commit ID"
// @Produce      json
// @Success      200  {object} []models.Thread
// @Failure      400  {object} httperror.HTTPError
// @Router       /articles/:articleID/versions/:versionID/history/:commitID/threads [get]
func (contr *ThreadController) GetCommitThreads(c *gin.Context) {
	aid := c.Param("articleID")
	cid := c.Param("commitID")

	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("articleID id %v is invalid", aid))
		return
	}

	if !gitUtils.IsCommitHash(cid) {
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("commit id %v is invalid", cid))
		return
	}

	threads, err := contr.CommitThreadService.GetCommitThreads(intAid, cid)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, fmt.Errorf("cannot find committhreads for this article"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, threads)
}

// GetCommitSelectionThreads godoc
// @Summary      Get all section threads for a commit
// @Description  Gets a list with all threads belonging to a specific commit of an article
// @Param		 article ID 		path	int64		true 	"Article ID"
// @Param		 commit ID 			path	int64		true 	"Commit ID"
// @Produce      json
// @Success      200  {object} []models.SelectionThread
// @Failure      400  {object} httperror.HTTPError
// @Router       /articles/:articleID/versions/:versionID/history/:commitID/selectionThreads [get]
func (contr *ThreadController) GetCommitSelectionThreads(c *gin.Context) {
	aid := c.Param("articleID")
	cid := c.Param("commitID")

	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("articleID id %v is invalid", aid))
		return
	}

	if !gitUtils.IsCommitHash(cid) {
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("commit id %v is invalid", cid))
		return
	}

	threads, err := contr.CommitSelectionThreadService.GetCommitSelectionThreads(cid, intAid)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, fmt.Errorf("cannot find committhreads for this article"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, threads)
}
