package thread

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/tests"
	"mainServer/tests/mocks/services"
	"net/http"
	"testing"
)

var threadServMock services.ThreadServiceMock
var commitThreadServMock services.CommitThreadServiceMock
var requestThreadServMock services.RequestThreadServiceMock
var contr controllers.ThreadController
var r *gin.Engine
var loggedInUser *string

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test in this file starts
func globalSetup() {
	// Setup test router, to test controller endpoints through http
	r = gin.Default()
	contr = controllers.ThreadController{}

	// Mock the authentication middleware, that sets the email of logged-in user
	r.Use(func(c *gin.Context) {
		if loggedInUser != nil {
			c.Set("Email", *loggedInUser)
		}
	})

	// routes
	r.POST("/articles/:articleID/thread/:threadType/id/:specificID/", func(c *gin.Context) {
		contr.CreateThread(c)
	})

	r.GET("/articles/:articleID/versions/:versionID/history/:commitID/threads", func(c *gin.Context) {
		contr.GetCommitThreads(c)
	})

	r.GET("/articles/:articleID/requests/:requestID/threads", func(c *gin.Context) {
		contr.GetRequestThreads(c)
	})

	r.POST("/comments/thread/:threadID", func(c *gin.Context) {
		contr.SaveComment(c)
	})
}

// localSetup should be called before each individual test
func localSetup() {
	//(Re)set controller mocks
	threadServMock = services.NewThreadServiceMock()
	requestThreadServMock = services.NewRequestThreadServiceMock()
	commitThreadServMock = services.NewCommitThreadServiceMock()

	contr.ThreadService = threadServMock
	contr.RequestThreadService = requestThreadServMock
	contr.CommitThreadService = commitThreadServMock
}

func TestGetCommitThreadFail(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitThreadsMock = func(aid int64, cid string) ([]models.Thread, error) {
		return []models.Thread{}, errors.New("cannot get commit threads")
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid string = "3"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/%s/threads", aid, vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusInternalServerError, nil)

	// check the service mock
	commitThreadServMock.Mock.AssertCalled(t, "GetCommitThreads", 1)
	commitThreadServMock.Mock.AssertCalledWith(t, "GetCommitThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
}

// TODO: getcommitthread
// bad param

func TestGetCommitThreadsSuccess(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitThreadsMock = func(aid int64, cid string) ([]models.Thread, error) {
		return exampleThreads, nil
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid string = "3"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/%s/threads", aid, vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusOK, exampleThreads)

	// check the service mock
	commitThreadServMock.Mock.AssertCalled(t, "GetCommitThreads", 1)
	commitThreadServMock.Mock.AssertCalledWith(t, "GetCommitThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
}

func TestGetRequestThreadFail(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetRequestThreadsMock = func(aid int64, rid int64) ([]models.Thread, error) {
		return []models.Thread{}, errors.New("cannot get request threads")
	}

	// set request url
	const aid int64 = 2
	const rid int64 = 3
	url := fmt.Sprintf("/articles/%d/requests/%d/threads", aid, rid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusInternalServerError, nil)

	// check the service mock
	requestThreadServMock.Mock.AssertCalled(t, "GetRequestThreads", 1)
	requestThreadServMock.Mock.AssertCalledWith(t, "GetRequestThreads", &map[string]interface{}{
		"aid": aid,
		"rid": rid,
	})
}

// TODO: getcommitthread
// bad params

func TestGetRequestThreadsSuccess(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetRequestThreadsMock = func(aid int64, rid int64) ([]models.Thread, error) {
		return exampleThreads, nil
	}

	// set request url
	const aid int64 = 2
	const rid int64 = 3
	url := fmt.Sprintf("/articles/%d/requests/%d/threads", aid, rid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusOK, exampleThreads)

	// check the service mock
	requestThreadServMock.Mock.AssertCalled(t, "GetRequestThreads", 1)
	requestThreadServMock.Mock.AssertCalledWith(t, "GetRequestThreads", &map[string]interface{}{
		"aid": aid,
		"rid": rid,
	})
}

// TODO: commitSelection
// success
// bad param
// 500

// TODO: saveCOmment

// TODo:

var exampleThreads = []models.Thread{
	{
		Id:         1,
		ArticleId:  2,
		SpecificId: "2",
		Comments: []entities.Comment{
			{
				Id:           1,
				AuthorId:     "pietje@gmail.com",
				ThreadId:     1,
				Content:      "Hey",
				CreationDate: "12345678",
			}, {
				Id:           2,
				AuthorId:     "puk@gmail.com",
				ThreadId:     1,
				Content:      "Hello!",
				CreationDate: "12345678",
			}},
	},
	{
		Id:         2,
		ArticleId:  2,
		SpecificId: "3",
		Comments: []entities.Comment{
			{
				Id:           1,
				AuthorId:     "pietje@gmail.com",
				ThreadId:     2,
				Content:      "Hey again",
				CreationDate: "12345678",
			}, {
				Id:           2,
				AuthorId:     "puk@gmail.com",
				ThreadId:     2,
				Content:      "Hello again!",
				CreationDate: "12345678",
			}},
	},
}
