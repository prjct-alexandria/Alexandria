package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services"
	"net/http"
)

type ThreadController struct {
	ThreadService       services.ThreadService
	CommitThreadService services.CommitThreadService
	//RequestThreadService services.RequestThreadService
	//ReviewThreadService  services.ReviewThreadService
	CommentService services.CommentService
}

// creates thread entity, and specific thread entity. returns both id's
func (contr *ThreadController) CreateThread(c *gin.Context) {
	var thread models.ThreadNoId
	err := c.BindJSON(&thread)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	aid := c.Param("articleID")
	cid := c.Param("commitID")
	threadType := c.Param("threadType")

	// save threat in the db
	tid, err := contr.ThreadService.StartThread(thread, aid, cid)

	// save first comment in the db
	coid, err := contr.CommentService.SaveComment(thread.Comment, tid)

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := int64(0), nil
	switch threadType {
	case "commit":
		id, err = contr.CommitThreadService.StartCommitThread(thread, tid)
		//case "request":
		//	id, err = contr.RequestThreadService.StartRequestThread(thread, tid)
		//case "review":
		//	id, err = contr.ReviewThreadService.StartReviewThread(thread, tid)
	}
	if err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// return tid and specific id
	ids := models.ReturnIds{
		ThreadId:   tid,
		CommentId:  coid,
		SpecificId: id,
	}

	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.IndentedJSON(http.StatusOK, ids)
}
