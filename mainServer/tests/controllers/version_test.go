package controllers

import (
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

// MainTest is a keyword function, this is run by the testing package
func MainTest(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestUpdateVersionSuccess(t *testing.T) {

	// define mock function behaviour
	updateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
		return nil
	}

	// Create controller with mock
	servMock := versionServiceMock{}
	contr := controllers.VersionController{Serv: servMock}

	// mock gin http request using built-in gin TestContext
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "articleID", Value: "42"},
		{Key: "versionID", Value: "123456"},
	}

	// act
	contr.UpdateVersion(c)

	// assert
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