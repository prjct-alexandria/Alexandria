package thread

import (
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
	"testing"
)

//var servMock services.ThreadServiceMock
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

	//// route
	//r.POST("/articles/:articleID/versions/:versionID", func(c *gin.Context) {
	//	contr.UpdateVersion(c)
	//})
}

// localSetup should be called before each individual test
func localSetup() {
	// (Re)set controller mocks
	//servMock = services.NewVersionServiceMock()
	//contr.Serv = servMock
}
