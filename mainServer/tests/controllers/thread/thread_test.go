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
var commitSelectionThreadServMock services.CommitSelectionThreadServiceMock
var commentServMock services.CommentServiceMock
var reviewThreadServMock services.ReviewThreadServiceMock
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
	r.POST("/articles/:articleID/thread/:threadType/id/:specificID", func(c *gin.Context) {
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

	r.GET("/articles/:articleID/versions/:versionID/history/:commitID/selectionThreads", func(c *gin.Context) {
		contr.GetCommitSelectionThreads(c)
	})
}

// localSetup should be called before each individual test
func localSetup() {
	//(Re)set controller mocks
	threadServMock = services.NewThreadServiceMock()
	requestThreadServMock = services.NewRequestThreadServiceMock()
	commitThreadServMock = services.NewCommitThreadServiceMock()
	commitSelectionThreadServMock = services.NewCommitSelectionThreadServiceMock()
	commentServMock = services.NewCommentServiceMock()
	reviewThreadServMock = services.NewReviewThreadServiceMock()

	contr.ThreadService = threadServMock
	contr.RequestThreadService = requestThreadServMock
	contr.CommitThreadService = commitThreadServMock
	contr.CommitSelectionThreadService = commitSelectionThreadServMock
	contr.CommentService = commentServMock
	contr.ReviewThreadService = reviewThreadServMock
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
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/%s/threads", aid, vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusOK, exampleThreads)

	// check the service mock
	commitThreadServMock.Mock.AssertCalled(t, "GetCommitThreads", 1)
	commitThreadServMock.Mock.AssertCalledWith(t, "GetCommitThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
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
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/%s/threads", aid, vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusInternalServerError, nil)

	// check the service mock
	commitThreadServMock.Mock.AssertCalled(t, "GetCommitThreads", 1)
	commitThreadServMock.Mock.AssertCalledWith(t, "GetCommitThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
}

func TestGetCommitThreadBadParam(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitThreadsMock = func(aid int64, cid string) ([]models.Thread, error) {
		return []models.Thread{}, errors.New("cannot get commit threads")
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/a/versions/%d/history/%s/threads", vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	commitThreadServMock.Mock.AssertCalled(t, "GetCommitThreads", 0)
}

func TestGetCommitThreadBadParam3(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitThreadsMock = func(aid int64, cid string) ([]models.Thread, error) {
		return []models.Thread{}, errors.New("cannot get commit threads")
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/1234/threads", aid, vid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	commitThreadServMock.Mock.AssertCalled(t, "GetCommitThreads", 0)
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

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusOK, exampleThreads)

	// check the service mock
	requestThreadServMock.Mock.AssertCalled(t, "GetRequestThreads", 1)
	requestThreadServMock.Mock.AssertCalledWith(t, "GetRequestThreads", &map[string]interface{}{
		"aid": aid,
		"rid": rid,
	})
}

func TestGetRequestThreadsBadParam(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetRequestThreadsMock = func(aid int64, rid int64) ([]models.Thread, error) {
		return exampleThreads, nil
	}

	// set request url
	const aid int64 = 2
	const rid int64 = 3
	url := fmt.Sprintf("/articles/a/requests/%d/threads", rid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	requestThreadServMock.Mock.AssertCalled(t, "GetRequestThreads", 0)
}

func TestGetRequestThreadsBadParam2(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetRequestThreadsMock = func(aid int64, rid int64) ([]models.Thread, error) {
		return exampleThreads, nil
	}

	// set request url
	const aid int64 = 2
	const rid int64 = 3
	url := fmt.Sprintf("/articles/%d/requests/a/threads", aid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	requestThreadServMock.Mock.AssertCalled(t, "GetRequestThreads", 0)
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

func TestGetCommitSelectionThreadsSuccess(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitSelectionThreadsMock = func(sid string, aid int64) ([]models.SelectionThread, error) {
		return exampleSelectionThreads, nil
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/%s/selectionThreads", aid, vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusOK, exampleSelectionThreads)

	// check the service mock
	commitSelectionThreadServMock.Mock.AssertCalled(t, "GetCommitSelectionThreads", 1)
	commitSelectionThreadServMock.Mock.AssertCalledWith(t, "GetCommitSelectionThreads", &map[string]interface{}{
		"sid": cid,
		"aid": aid,
	})
}

func TestGetCommitSelectionThreadFail(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitSelectionThreadsMock = func(sid string, aid int64) ([]models.SelectionThread, error) {
		return nil, errors.New("fake error for testing")
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/%s/selectionThreads", aid, vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusInternalServerError, nil)

	// check the service mock
	commitSelectionThreadServMock.Mock.AssertCalled(t, "GetCommitSelectionThreads", 1)
	commitSelectionThreadServMock.Mock.AssertCalledWith(t, "GetCommitSelectionThreads", &map[string]interface{}{
		"sid": cid,
		"aid": aid,
	})
}

func TestGetCommitSelectionThreadBadParam(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitSelectionThreadsMock = func(sid string, aid int64) ([]models.SelectionThread, error) {
		return exampleSelectionThreads, nil
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/a/versions/%d/history/%s/selectionThreads", vid, cid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	commitSelectionThreadServMock.Mock.AssertCalled(t, "GetCommitSelectionThreads", 0)
}

func TestGetCommitSelectionThreadBadParam3(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.GetCommitSelectionThreadsMock = func(sid string, aid int64) ([]models.SelectionThread, error) {
		return exampleSelectionThreads, nil
	}

	// set request url
	const aid int64 = 2
	const vid int64 = 1
	const cid = "0123456789012345678901234567890123456789"
	url := fmt.Sprintf("/articles/%d/versions/%d/history/a/selectionThreads", aid, vid)

	tests.TestEndpoint(t, r, "GET", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	commitSelectionThreadServMock.Mock.AssertCalled(t, "GetCommitSelectionThreads", 0)
}

func TestSaveCommentSuccess(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 42, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	const tid int64 = 2
	url := fmt.Sprintf("/comments/thread/%d", tid)

	tests.TestEndpoint(t, r, "POST", url, exampleComment, http.StatusOK, nil)

	// check the service mock
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	commentServMock.Mock.AssertCalledWith(t, "SaveComment", &map[string]interface{}{
		"comment":    exampleComment,
		"tid":        tid,
		"loggedInAs": email,
	})
}

func TestSaveCommentNotLoggedIn(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 42, nil
	}

	// mock the authentication and set the logged-in user
	loggedInUser = nil

	// set request url
	const tid int64 = 2
	url := fmt.Sprintf("/comments/thread/%d", tid)

	tests.TestEndpoint(t, r, "POST", url, exampleComment, http.StatusForbidden, nil)

	// check the service mock
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
}

func TestSaveCommentBadForm(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 42, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// make struct for comment that only has one field
	badComment := struct {
		Content string
	}{
		Content: "some text",
	}

	// set request url
	const tid int64 = 2
	url := fmt.Sprintf("/comments/thread/%d", tid)

	tests.TestEndpoint(t, r, "POST", url, badComment, http.StatusBadRequest, nil)

	// check the service mock
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
}

func TestSaveCommentBadParam(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 42, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	const tid int64 = 2
	url := fmt.Sprintf("/comments/thread/a")

	tests.TestEndpoint(t, r, "POST", url, exampleComment, http.StatusBadRequest, nil)

	// check the service mock
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
}

func TestSaveCommentInternalError(t *testing.T) {
	localSetup()

	//define service mock behaviour
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return -1, errors.New("fake error for testing")
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	const tid int64 = 2
	url := fmt.Sprintf("/comments/thread/%d", tid)

	tests.TestEndpoint(t, r, "POST", url, exampleComment, http.StatusInternalServerError, nil)

	// check the service mock
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	commentServMock.Mock.AssertCalledWith(t, "SaveComment", &map[string]interface{}{
		"comment":    exampleComment,
		"tid":        tid,
		"loggedInAs": email,
	})
}

func TestCreateCommitThread(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid string = "0123456789012345678901234567890123456789"
	const tid int64 = 42

	//define service mock behaviour
	services.StartThreadMock = func(thread models.Thread, aid int64) (int64, error) {
		return 42, nil
	}
	services.StartCommitThreadMock = func(sid string, tid int64) (int64, error) {
		return 43, nil
	}
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 44, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/commit/id/%s", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleThreads[0], http.StatusOK, exampleReturnThreadIds)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 1)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	commitThreadServMock.Mock.AssertCalled(t, "StartCommitThread", 1)
	commitThreadServMock.Mock.AssertCalledWith(t, "StartCommitThread", &map[string]interface{}{
		"cid": sid,
		"tid": tid,
	})
}

func TestCreateCommitSelectionThread(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid string = "0123456789012345678901234567890123456789"
	const tid int64 = 42

	//define service mock behaviour
	services.StartThreadMock = func(thread models.Thread, aid int64) (int64, error) {
		return 42, nil
	}
	services.StartCommitSelectionThreadMock = func(cid string, tid int64, selection string) (int64, error) {
		return 43, nil
	}
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 44, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/commitSelection/id/%s", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusOK, exampleReturnThreadIds)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 1)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	commitSelectionThreadServMock.Mock.AssertCalled(t, "StartCommitSelectionThread", 1)
	commitSelectionThreadServMock.Mock.AssertCalledWith(t, "StartCommitSelectionThread", &map[string]interface{}{
		"cid":     sid,
		"tid":     tid,
		"section": exampleSelectionThreads[0].Selection,
	})
}

func TestCreateRequestThread(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	//define service mock behaviour
	services.StartThreadMock = func(thread models.Thread, aid int64) (int64, error) {
		return 42, nil
	}
	services.StartRequestThreadMock = func(rid int64, tid int64, loggedInAs string) (int64, error) {
		return 43, nil
	}
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 44, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/request/id/%d", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusOK, exampleReturnThreadIds)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 1)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	requestThreadServMock.Mock.AssertCalled(t, "StartRequestThread", 1)
	requestThreadServMock.Mock.AssertCalledWith(t, "StartRequestThread", &map[string]interface{}{
		"rid":        sid,
		"tid":        tid,
		"loggedInAs": email,
	})
}

func TestCreateReviewThread(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	//define service mock behaviour
	services.StartThreadMock = func(thread models.Thread, aid int64) (int64, error) {
		return 42, nil
	}
	services.StartReviewThreadMock = func(rid int64, tid int64) (int64, error) {
		return 43, nil
	}
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return 44, nil
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/review/id/%d", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusOK, exampleReturnThreadIds)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 1)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	reviewThreadServMock.Mock.AssertCalled(t, "StartReviewThread", 1)
	reviewThreadServMock.Mock.AssertCalledWith(t, "StartReviewThread", &map[string]interface{}{
		"rid": sid,
		"tid": tid,
	})
}

func TestCreateThreadNotLoggedIn(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	// mock the authentication and set the logged-in user
	loggedInUser = nil

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/review/id/%d", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusForbidden, nil)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 0)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
	reviewThreadServMock.Mock.AssertCalled(t, "StartReviewThread", 0)
}

func TestCreateThreadBadJson(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/review/id/%d", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, "some literal comment text", http.StatusBadRequest, nil)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 0)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
	reviewThreadServMock.Mock.AssertCalled(t, "StartReviewThread", 0)
}

func TestCreateThreadBadParameter(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/a/thread/review/id/%d", sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusBadRequest, nil)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 0)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
	reviewThreadServMock.Mock.AssertCalled(t, "StartReviewThread", 0)
}

func TestCreateThreadThreadError(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	services.StartThreadMock = func(thread models.Thread, aid int64) (int64, error) {
		return -1, errors.New("fake error for testing")
	}

	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/review/id/%d", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusInternalServerError, nil)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 1)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 0)
	reviewThreadServMock.Mock.AssertCalled(t, "StartReviewThread", 0)
}

func TestCreateThreadCommentError(t *testing.T) {
	localSetup()

	const aid int64 = 2
	const sid int64 = 9
	const tid int64 = 42

	//define service mock behaviour
	services.StartThreadMock = func(thread models.Thread, aid int64) (int64, error) {
		return 42, nil
	}
	services.StartReviewThreadMock = func(rid int64, tid int64) (int64, error) {
		return 43, nil
	}
	services.SaveCommentMock = func(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
		return -1, errors.New("fake error for testing")
	}
	// mock the authentication and set the logged-in user
	email := "john@mail.com"
	loggedInUser = &email

	// set request url
	url := fmt.Sprintf("/articles/%d/thread/review/id/%d", aid, sid)

	tests.TestEndpoint(t, r, "POST", url, exampleSelectionThreads[0], http.StatusInternalServerError, nil)

	// check the service mocks
	threadServMock.Mock.AssertCalled(t, "StartThread", 1)
	commentServMock.Mock.AssertCalled(t, "SaveComment", 1)
	reviewThreadServMock.Mock.AssertCalled(t, "StartReviewThread", 0)
}

// success
// not logged in
// bad aid
// 500 thread
// 500 comment
// bad cid
// bad selection
// 500 ___threadError
// bad threadType

var exampleComment = entities.Comment{
	Id:           1,
	AuthorId:     "john@mail.com",
	ThreadId:     2,
	Content:      "I like this",
	CreationDate: "2022-06-26T22:20:17.485Z",
}

var exampleReturnThreadIds = models.ReturnThreadIds{
	Id:        43,
	ThreadId:  42,
	CommentId: 44,
}

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

var exampleSelectionThreads = []models.SelectionThread{
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
		Selection: "some quote",
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
		Selection: "some other quote",
	},
}
