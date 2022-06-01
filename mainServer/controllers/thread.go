package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services"
	"net/http"
	"strconv"
)

type ThreadController struct {
	ThreadService        services.ThreadService
	CommitThreadService  services.CommitThreadService
	RequestThreadService services.RequestThreadService
	ReviewThreadService  services.ReviewThreadService
	CommentService       services.CommentService
}

// CreateThread godoc
// @Summary      Creates thread entity
// @Description  Creates thread entity, and specific thread entity. Returns id's of thread, specific thread and comment
// @Accept		 json
// @Param		 thread 		body	models.Thread		true 	"Thread"
// @Produce      json
// @Success      200  {object} models.ReturnIds
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

	// save threat in the db
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
