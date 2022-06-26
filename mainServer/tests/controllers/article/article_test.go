package article

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

var articleServMock services.ArticleServiceMock
var contr *controllers.ArticleController
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
	r.POST("/articles", func(c *gin.Context) {
		contr.CreateArticle(c)
	})
	r.GET("/articles/:articleID/mainVersion", func(c *gin.Context) {
		contr.GetMainVersion(c)
	})
	r.GET("/articles", func(c *gin.Context) {
		contr.ArticleList(c)
	})

}

// localSetup should be called before each individual test
func localSetup() {
	// Make a clean controller with clean mocks
	articleServMock = services.NewArticleServiceMock()
	contrVal := controllers.NewArticleController(articleServMock)
	contr = &contrVal
}

func TestCreateArticleSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateArticleMock = func(title string, owners []string, loggedInAs string) (models.Version, error) {
		return models.Version{Title: title, Owners: owners}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// define expected http response body
	expectedVersion := models.Version{Title: exampleArticleCreation.Title, Owners: exampleArticleCreation.Owners}

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles", exampleArticleCreation,
		http.StatusOK, expectedVersion)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "CreateArticle", 1)
	articleServMock.Mock.AssertCalledWith(t, "CreateArticle", &map[string]interface{}{
		"title":      exampleArticleCreation.Title,
		"owners":     exampleArticleCreation.Owners,
		"loggedInAs": user,
	})
}

func TestCreateArticleNotLoggedIn(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateArticleMock = func(title string, owners []string, loggedInAs string) (models.Version, error) {
		return models.Version{Title: title, Owners: owners}, nil
	}

	// fake auth middleware and set the logged-in user
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles", exampleArticleCreation,
		http.StatusForbidden, nil)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "CreateArticle", 0)
}

func TestCreateArticleMissingFields(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateArticleMock = func(title string, owners []string, loggedInAs string) (models.Version, error) {
		return models.Version{Title: title, Owners: owners}, nil
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// Create incomplete article creation form
	badArticleCreation := struct {
		Title string `json:"title"`
		// note: left out owners here
	}{
		Title: exampleArticleCreation.Title,
	}

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles", badArticleCreation,
		http.StatusBadRequest, nil)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "CreateArticle", 0)
}

func TestCreateArticleInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.CreateArticleMock = func(title string, owners []string, loggedInAs string) (models.Version, error) {
		return models.Version{}, errors.New("fake error for testing")
	}

	// fake auth middleware and set the logged-in user
	user := "john@mail.com"
	loggedInUser = &user

	// execute and test the http request
	tests.TestEndpoint(t, r, "POST", "/articles", exampleArticleCreation,
		http.StatusInternalServerError, nil)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "CreateArticle", 1)
	articleServMock.Mock.AssertCalledWith(t, "CreateArticle", &map[string]interface{}{
		"title":      exampleArticleCreation.Title,
		"owners":     exampleArticleCreation.Owners,
		"loggedInAs": user,
	})
}

func TestGetMainVersionSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetMainVersionMock = func(article int64) (int64, error) {
		return 42, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/36/mainVersion", exampleArticleCreation,
		http.StatusOK, "42")

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "GetMainVersion", 1)
	articleServMock.Mock.AssertCalledWith(t, "GetMainVersion", &map[string]interface{}{
		"article": int64(36),
	})
}

func TestGetMainVersionBadId(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetMainVersionMock = func(article int64) (int64, error) {
		return 42, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/3a6/mainVersion", exampleArticleCreation,
		http.StatusBadRequest, nil)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "GetMainVersion", 0)
}

func TestGetMainVersionInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetMainVersionMock = func(article int64) (int64, error) {
		return -1, errors.New("fake error for testing")
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/36/mainVersion", exampleArticleCreation,
		http.StatusInternalServerError, nil)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "GetMainVersion", 1)
	articleServMock.Mock.AssertCalledWith(t, "GetMainVersion", &map[string]interface{}{
		"article": int64(36),
	})
}

func TestArticleListSuccess(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetArticleListMock = func() ([]models.ArticleListElement, error) {
		return exampleArticleList, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles", exampleArticleCreation,
		http.StatusOK, exampleArticleList)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "GetArticleList", 1)
}

func TestArticleListInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetArticleListMock = func() ([]models.ArticleListElement, error) {
		return nil, errors.New("fake error for testing")
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles", exampleArticleCreation,
		http.StatusInternalServerError, nil)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "GetArticleList", 1)
}

var exampleArticleCreation = models.ArticleCreationForm{
	Title:  "Example Article",
	Owners: []string{"john@mail.com", "jane@mail.com"},
}

var exampleArticleList = []models.ArticleListElement{
	{
		Id:            1,
		MainVersionId: 1,
		Title:         "Awesome article",
		Owners:        []string{"jane@mail.com", "joe@mail.com"},
	},
	{
		Id:            2,
		MainVersionId: 72,
		Title:         "Awesome article about alexandria",
		Owners:        []string{"jane@mail.com", "jill@mail.com"},
	},
}
