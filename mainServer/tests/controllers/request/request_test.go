package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
	"mainServer/models"
	"mainServer/tests"
	"mainServer/tests/mocks/services"
	"net/http"
	"testing"
)

var requestServMock services.RequestServiceMock
var contr *controllers.RequestController
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

	// Mock the authentication middleware, that sets the email of logged-in user
	r.Use(func(c *gin.Context) {
		if loggedInUser != nil {
			c.Set("Email", *loggedInUser)
		}
	})

	// route
	r.POST("/articles/:articleID/requests", func(c *gin.Context) {
		contr.CreateRequest(c)
	})
	r.PUT("/articles/:articleID/requests/:requestID/reject", func(c *gin.Context) {
		contr.RejectRequest(c)
	})
	r.PUT("/articles/:articleID/requests/:requestID/accept", func(c *gin.Context) {
		contr.AcceptRequest(c)
	})
}

// localSetup should be called before each individual test
func localSetup() {
	// Make a clean controller with clean mocks
	requestServMock = services.NewRequestServiceMock()
	contrVal := controllers.RequestController{Serv: requestServMock}
	contr = &contrVal
}

func TestCreateRequestSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return exampleRequest, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// define expected http response body
	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles/2/requests", exampleRequestCreation,
		http.StatusOK, exampleRequest)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "CreateRequest", 1)
	requestServMock.Mock.AssertCalledWith(t, "CreateRequest", &map[string]interface{}{
		"article":       int64(2),
		"sourceVersion": exampleRequestCreation.SourceVersionID,
		"targetVersion": exampleRequestCreation.TargetVersionID,
		"loggedInAs":    user,
	})
}

func TestCreateRequestNotLoggedIn(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return exampleRequest, nil
	}

	// fake auth middleware and set the logged-in user
	loggedInUser = nil

	// define expected http response body
	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles/2/requests", exampleRequestCreation,
		http.StatusForbidden, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "CreateRequest", 0)
}

func TestCreateRequestBadForm(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return exampleRequest, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// make request creation form with missing fields
	badForm := struct {
		SourceVersionID int64
	}{
		SourceVersionID: 1,
	}

	// define expected http response body
	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles/2/requests", badForm,
		http.StatusBadRequest, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "CreateRequest", 0)
}

func TestCreateRequestBadParameter(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return models.Request{}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// define expected http response body
	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles/a/requests", exampleRequestCreation,
		http.StatusBadRequest, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "CreateRequest", 0)
}

func TestCreateRequestInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return models.Request{}, errors.New("fake error for testing")
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// define expected http response body
	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles/2/requests", exampleRequestCreation,
		http.StatusInternalServerError, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "CreateRequest", 1)
	requestServMock.Mock.AssertCalledWith(t, "CreateRequest", &map[string]interface{}{
		"article":       int64(2),
		"sourceVersion": exampleRequestCreation.SourceVersionID,
		"targetVersion": exampleRequestCreation.TargetVersionID,
		"loggedInAs":    user,
	})
}

func TestAcceptRequestSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.AcceptRequestMock = func(request int64, loggedInAs string) error {
		return nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/2/requests/1/accept", nil,
		http.StatusOK, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "AcceptRequest", 1)
	requestServMock.Mock.AssertCalledWith(t, "AcceptRequest", &map[string]interface{}{
		"request":    int64(1),
		"loggedInAs": user,
	})
}

func TestAcceptRequestNotLoggedIn(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.AcceptRequestMock = func(request int64, loggedInAs string) error {
		return nil
	}

	// fake auth middleware and set the logged-in user
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/2/requests/1/accept", nil,
		http.StatusForbidden, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "AcceptRequest", 0)
}

func TestAcceptRequestBadParameter(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return models.Request{}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/a/requests/1/accept", nil,
		http.StatusBadRequest, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "AcceptRequest", 0)
}

func TestAcceptRequestBadParameter2(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return models.Request{}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/1/requests/a/accept", nil,
		http.StatusBadRequest, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "AcceptRequest", 0)
}

func TestAcceptRequestInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.AcceptRequestMock = func(request int64, loggedInAs string) error {
		return errors.New("fake error for testing")
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/2/requests/1/accept", nil,
		http.StatusInternalServerError, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "AcceptRequest", 1)
	requestServMock.Mock.AssertCalledWith(t, "AcceptRequest", &map[string]interface{}{
		"request":    int64(1),
		"loggedInAs": user,
	})
}

func TestRejectRequestSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.RejectRequestMock = func(request int64, loggedInAs string) error {
		return nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/2/requests/1/reject", nil,
		http.StatusOK, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "RejectRequest", 1)
	requestServMock.Mock.AssertCalledWith(t, "RejectRequest", &map[string]interface{}{
		"request":    int64(1),
		"loggedInAs": user,
	})
}

func TestRejectRequestNotLoggedIn(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.RejectRequestMock = func(request int64, loggedInAs string) error {
		return nil
	}

	// fake auth middleware and set the logged-in user
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/2/requests/1/reject", nil,
		http.StatusForbidden, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "RejectRequest", 0)
}

func TestRejectRequestBadParameter(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return models.Request{}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/a/requests/1/reject", nil,
		http.StatusBadRequest, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "RejectRequest", 0)
}

func TestRejectRequestBadParameter2(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateRequestMock = func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
		return models.Request{}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/1/requests/a/reject", nil,
		http.StatusBadRequest, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "RejectRequest", 0)
}

func TestRejectRequestInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.RejectRequestMock = func(request int64, loggedInAs string) error {
		return errors.New("fake error for testing")
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "PUT", "/articles/2/requests/1/reject", nil,
		http.StatusInternalServerError, nil)

	// check the service mock
	requestServMock.Mock.AssertCalled(t, "RejectRequest", 1)
	requestServMock.Mock.AssertCalledWith(t, "RejectRequest", &map[string]interface{}{
		"request":    int64(1),
		"loggedInAs": user,
	})
}

// TODO: getRequest
// success
// bad parameters
// internal error

// TODO: getRequest
// success
// bad parameters
// internal error

var exampleRequestCreation = models.RequestCreationForm{
	SourceVersionID: 1,
	TargetVersionID: 2,
}

var exampleRequest = models.Request{
	RequestID:       1,
	ArticleID:       2,
	SourceVersionID: exampleRequestCreation.SourceVersionID,
	SourceHistoryID: "0123456789012345678901234567890123456789",
	TargetVersionID: exampleRequestCreation.TargetVersionID,
	TargetHistoryID: "0123456789012345678901234567890123456789",
	Status:          "draft",
	Conflicted:      false,
}
