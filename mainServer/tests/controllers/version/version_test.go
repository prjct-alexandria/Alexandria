package version

import (
	"bytes"
	"errors"
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

// globalSetup should be called once, before any test in this file starts
func globalSetup() {
	// Setup test router, to test controller endpoints through http
	r = gin.Default()
	contr = controllers.VersionController{}

	// route
	r.POST("/articles/:articleID/versions/:versionID", func(c *gin.Context) {
		contr.UpdateVersion(c)
	})
}

// localSetup should be called before each individual test
func localSetup() {
	// (Re)set controller mocks
	servMock = services.NewVersionServiceMock()
	contr.Serv = servMock
}

func TestUpdateVersionSuccess(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {
		return nil
	}

	// set request url
	const article int64 = 42
	const version int64 = 123456
	url := fmt.Sprintf("/articles/%d/versions/%d", article, version)

	// create request
	req, err := fileUploadHelper(url)
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
	if !(*servMock.Called)["UpdateVersion"] {
		t.Errorf("Expected UpdateVersion to be called")
	}
	a := (*servMock.Params)["UpdateVersion"]["article"]
	v := (*servMock.Params)["UpdateVersion"]["version"]
	if a != article || v != version {
		t.Errorf("Expected different function params, got article=%v and version=%v", a, v)
	}
}

func TestUpdateVersionFail(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {
		return errors.New("oh no, the version coulnd't be updated")
	}

	// set request url
	const article int64 = 42
	const version int64 = 123456
	url := fmt.Sprintf("/articles/%d/versions/%d", article, version)

	// create request
	req, err := fileUploadHelper(url)
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
	if !(*servMock.Called)["UpdateVersion"] {
		t.Errorf("Expected UpdateVersion to be called")
	}
	a := (*servMock.Params)["UpdateVersion"]["article"]
	v := (*servMock.Params)["UpdateVersion"]["version"]
	if a != article || v != version {
		t.Errorf("Expected different function params, got article=%v and version=%v", a, v)
	}
}

func TestUpdateVersionNoFile(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {
		return nil
	}

	// set request url
	const article = "42"
	const version = "123456"
	url := fmt.Sprintf("/articles/%s/versions/%s", article, version)

	// create request
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Error(err)
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
	if (*servMock.Called)["UpdateVersion"] {
		t.Errorf("Expected UpdateVersion to not be called")
	}
}

// fileUploaderHelper creates a http request with a file in the form data
func fileUploadHelper(url string) (*http.Request, error) {
	// Use the multipart format to attach a file
	var b bytes.Buffer
	m := multipart.NewWriter(&b)
	field, err := m.CreateFormFile("file", "helloworld.txt")
	if err != nil {
		return &http.Request{}, err
	}

	// To add some file contents, write "Hello World" in ASCII,
	// directly as a byte array for simplicity, instead of reading from an external file
	_, err = field.Write([]byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100})
	if err != nil {
		return &http.Request{}, err
	}
	err = m.Close()
	if err != nil {
		return &http.Request{}, err
	}

	// create request
	req, err := http.NewRequest(http.MethodPost, url, &b)
	req.Header.Set("Content-Type", m.FormDataContentType())
	if err != nil {
		return nil, err
	}
	return req, nil
}
