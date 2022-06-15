package thread

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/controllers"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/tests/mocks/services"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var threadServMock services.ThreadServiceMock
var commitThreadServMock services.CommitThreadServiceMock
var requestThreadServMock services.RequestThreadServiceMock
var contr controllers.ThreadController
var r *gin.Engine

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

	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
		return
	}

	// perform the request, with a response recorder
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// check the response
	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	// check the service mock
	if !(*commitThreadServMock.Called)["GetCommitThreads"] {
		t.Errorf("Expected GetRequestThreads to be called")
	}
	a := (*commitThreadServMock.Params)["GetCommitThreads"]["aid"]
	c := (*commitThreadServMock.Params)["GetCommitThreads"]["cid"]
	if a != aid || c != cid {
		t.Errorf("Expected different function params, got article=%v and commit=%v", a, c)
	}
}

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

	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
		return
	}

	// perform the request with a response recorder
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// check the response
	b, err := ioutil.ReadAll(w.Body)
	if w.Code != 200 {
		t.Error(w.Code, string(b))
	}
	var threads []models.Thread
	err = json.Unmarshal(b, &threads)
	if !reflect.DeepEqual(threads, exampleThreads) {
		t.Errorf("Expected expected=%v but got actual=%v", exampleThreads, threads)
	}

	// check the service mock
	if !(*commitThreadServMock.Called)["GetCommitThreads"] {
		t.Errorf("Expected GetCommitThreads to be called")
	}
	a := (*commitThreadServMock.Params)["GetCommitThreads"]["aid"]
	c := (*commitThreadServMock.Params)["GetCommitThreads"]["cid"]
	if a != aid || c != cid {
		t.Errorf("Expected different function params, got article=%v and commit=%v", a, c)
	}
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

	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
		return
	}

	// perform the request, with a response recorder
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// check the response
	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	// check the service mock
	if !(*requestThreadServMock.Called)["GetRequestThreads"] {
		t.Errorf("Expected GetRequestThreads to be called")
	}
	a := (*requestThreadServMock.Params)["GetRequestThreads"]["aid"]
	r := (*requestThreadServMock.Params)["GetRequestThreads"]["rid"]
	if a != aid || r != rid {
		t.Errorf("Expected different function params, got article=%v and request=%v", a, r)
	}
}

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

	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
		return
	}

	// perform the request with a response recorder
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// check the response
	b, err := ioutil.ReadAll(w.Body)
	if w.Code != 200 {
		t.Error(w.Code, string(b))
	}
	var threads []models.Thread
	err = json.Unmarshal(b, &threads)
	if !reflect.DeepEqual(threads, exampleThreads) {
		t.Errorf("Expected expected=%v but got actual=%v", exampleThreads, threads)
	}

	// check the service mock
	if !(*requestThreadServMock.Called)["GetRequestThreads"] {
		t.Errorf("Expected GetRequestThreads to be called")
	}
	a := (*requestThreadServMock.Params)["GetRequestThreads"]["aid"]
	r := (*requestThreadServMock.Params)["GetRequestThreads"]["rid"]
	if a != aid || r != rid {
		t.Errorf("Expected different function params, got article=%v and request=%v", a, r)
	}
}

var exampleThreads = []models.Thread{
	{
		Id:         1,
		ArticleId:  2,
		SpecificId: "2",
		Comment: []entities.Comment{
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
		Comment: []entities.Comment{
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
