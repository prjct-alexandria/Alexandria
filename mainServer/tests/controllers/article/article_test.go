package article

import (
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
	"mainServer/models"
	"mainServer/tests"
	"mainServer/tests/mocks/services"
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

	// route
	r.POST("/articles", func(c *gin.Context) {
		if loggedInUser != nil {
			c.Set("Email", *loggedInUser)
		}
		contr.CreateArticle(c)
	})
}

// localSetup should be called before each individual test
func localSetup() {
	// Make a clean controller with clean mocks
	articleServMock = services.NewArticleServiceMock()
	contrVal := controllers.NewArticleController(articleServMock)
	contr = &contrVal
}

// TODO: createArticle
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
		200, expectedVersion)

	// check the service mock
	articleServMock.Mock.AssertCalled(t, "CreateArticle", 1)
	articleServMock.Mock.AssertCalledWith(t, "CreateArticle", &map[string]interface{}{
		"title":      exampleArticleCreation.Title,
		"owners":     exampleArticleCreation.Owners,
		"loggedInAs": user,
	})
}

// not logged in
// missing json fields
// 500

// TODO: get main version
// success
// bad id
// 500

// TODO: articlelist
// success
// 500

var exampleArticleCreation = models.ArticleCreationForm{
	Title:  "Example Article",
	Owners: []string{"john@mail.com", "jane@mail.com"},
}
