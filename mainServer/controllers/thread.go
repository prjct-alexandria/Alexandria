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

// CreateThread creates thread entity, and specific thread entity.
// returns id's of thread, specific thread and comment
func (contr *ThreadController) CreateThread(c *gin.Context) {
	var thread models.ThreadNoId
	err := c.BindJSON(&thread)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	aid := c.Param("articleID")
	sid := c.Param("specificID")
	threadType := c.Param("threadType")

	// save threat in the db
	tid, err := contr.ThreadService.StartThread(thread, aid, sid)
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

	id, err := int64(0), nil
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
	ids := models.ReturnIds{
		ThreadId:  tid,
		CommentId: coid,
		Id:        id,
	}

	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.IndentedJSON(http.StatusOK, ids)
}

func (contr *ThreadController) SaveComment(c *gin.Context) {
	var comment models.CommentNoId
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
		return
	}
	id, err := contr.CommentService.SaveComment(comment, intTid)

	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.IndentedJSON(http.StatusOK, id)
}
