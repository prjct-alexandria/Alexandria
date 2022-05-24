package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/controllers"
	"mainServer/tests/mocks/services"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var servMock services.VersionServiceMock
var contr controllers.VersionController
var r *gin.Engine

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test starts
func globalSetup() {
	// Setup test router, to test controller endpoints through http
	r = gin.Default()
}

// localSetup should be called before each individual test
func localSetup() {
	// Create controller with mocks
	contr = controllers.VersionController{}
	contr.Serv = services.VersionServiceMock{}

	// route
	r.POST("/articles/:articleID/versions/:versionID", func(c *gin.Context) {
		contr.UpdateVersion(c)
	})
}

func TestUpdateVersionSuccess(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
		return nil
	}

	// set request url
	const article = 42
	const version = 123456
	url := fmt.Sprintf("/articles/%d/versions/%d", article, version)

	// set request file contents as form data
	var b bytes.Buffer
	m := multipart.NewWriter(&b)
	field, err := m.CreateFormFile("file", "helloworld.txt")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = field.Write([]byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100})
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// create request
	req, err := http.NewRequest(http.MethodPost, url, &b)
	req.Header.Set("Content-Type", m.FormDataContentType())
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
		return
	}

	// perform the request, with a response recorder
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// check the response
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	// check the service mock
	if !servMock.UpdateVersionCalled {
		t.Errorf("Expected UpdateVersion to be called")
	}
	a := servMock.UpdateVersionParams["article"]
	v := servMock.UpdateVersionParams["version"]
	if a != article || v != version {
		t.Errorf("Expected different function params, got article=%v and version=%v", a, v)
	}
}
