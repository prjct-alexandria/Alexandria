package controllers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/controllers"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var servMock versionServiceMock
var contr controllers.VersionController
var r gin.Engine

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test starts
func globalSetup() {
	// Create controller, mocks are added later during localSetup
	contr = controllers.VersionController{}

	// Setup test router, to test controller endpoints through http
	r := gin.Default()
	r.POST("/articles/:articleID/versions/:versionID", contr.UpdateVersion)
}

// localSetup should be called before each individual test
func localSetup() {
	// initialize/reset mocks
	contr.Serv = versionServiceMock{}
}

func TestUpdateVersionSuccess(t *testing.T) {
	localSetup()

	// define service mock behaviour
	updateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
		return nil
	}

	// Create the mock request
	req, err := http.NewRequest(http.MethodPost, "/articles/:articleID/versions/:versionID", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Perform the request, with a response recorder
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the response
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	// Check the service mock
	if !servMock.UpdateVersionCalled {
		t.Errorf("Expected UpdateVersion to be called")
	}
	a := servMock.UpdateVersionParams["article"]
	v := servMock.UpdateVersionParams["version"]
	if a != 42 || v != 123456 {
		t.Errorf("Expected different function params, got article=%v and version=%v", a, v)
	}
}

// Declare mock functions without defining the body,
// so each test can specify the behaviour differently
var updateVersionMock func(c *gin.Context, file *multipart.FileHeader, article string, version string) error

// Mock class using the variable mock functions
type versionServiceMock struct {
	// mock tracks what functions were called
	UpdateVersionCalled bool
	UpdateVersionParams map[string]interface{}
}

func (m versionServiceMock) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
	m.UpdateVersionCalled = true
	m.UpdateVersionParams["article"] = article
	m.UpdateVersionParams["version"] = version
	return updateVersionMock(c, file, article, version)
}
