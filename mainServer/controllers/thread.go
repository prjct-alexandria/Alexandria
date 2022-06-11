package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/services/interfaces"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type ThreadController struct {
	ThreadService        interfaces.ThreadService
	CommitThreadService  interfaces.CommitThreadService
	RequestThreadService interfaces.RequestThreadService
	ReviewThreadService  interfaces.ReviewThreadService
	CommentService       interfaces.CommentService
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

	intSid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// save thread in the db
	tid, err := contr.ThreadService.StartThread(thread, intAid, intSid)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// save first comment in the db
	coid, err := contr.CommentService.SaveComment(thread.Comment[0], tid)

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	var id int64
	switch threadType {
	case "commit":
		id, err = contr.CommitThreadService.StartCommitThread(thread, tid)
	case "request":
		id, err = contr.RequestThreadService.StartRequestThread(thread, tid)
	case "review":
		id, err = contr.ReviewThreadService.StartReviewThread(thread, tid)
	default:
		id, err = -1, errors.New("invalid thread type")
	}

	if err != nil {
		c.Status(http.StatusBadRequest)
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
	c.IndentedJSON(http.StatusOK, ids)
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
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := contr.CommentService.SaveComment(comment, intTid)
	if err != nil {
		fmt.Println(err)
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
// @Success      200  {object} models.Thread
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
	if err != nil || threads == nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("cannot find requestthreads for this request"))
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
// @Success      200  {object} models.Thread
// @Failure      400  {object} httperror.HTTPError
// @Router       /articles/:articleID/versions/:versionID/history/:commitID/threads [get]
func (contr *ThreadController) GetCommitThreads(c *gin.Context) {
	aid := c.Param("articleID")
	cid := c.Param("commitID")

	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	intCid, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	threads, err := contr.CommitThreadService.GetCommitThreads(intAid, intCid)
	if err != nil || threads == nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("cannot find committhreads for this article"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, threads)
}
