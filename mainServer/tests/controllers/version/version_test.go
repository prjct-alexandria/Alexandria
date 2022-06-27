package version

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/controllers"
	"mainServer/models"
	"mainServer/tests"
	"mainServer/tests/mocks/services"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var servMock services.VersionServiceMock
var contr controllers.VersionController
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
	contr = controllers.VersionController{}

	// Mock the authentication middleware, that sets the email of logged-in user
	r.Use(func(c *gin.Context) {
		if loggedInUser != nil {
			c.Set("Email", *loggedInUser)
		}
	})

	// route
	r.POST("/articles/:articleID/versions/:versionID", func(c *gin.Context) {
		contr.UpdateVersion(c)
	})
	r.GET("/articles/:articleID/versions/:versionID", func(c *gin.Context) {
		contr.GetVersion(c)
	})
	r.GET("/articles/:articleID", func(c *gin.Context) {
		contr.ListVersions(c)
	})
	r.POST("/articles/:articleID/versions", func(c *gin.Context) {
		contr.CreateVersionFrom(c)
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
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
		return nil
	}

	// set the logged-in user
	email := "zoo@mail.com"
	loggedInUser = &email

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
	servMock.Mock.AssertCalled(t, "UpdateVersion", 1)
	// check the service mock params manually, only checks the non-pointer parameters
	calledWith := servMock.Mock.Params["UpdateVersion"]
	if calledWith["article"] != article ||
		calledWith["version"] != version ||
		calledWith["loggedInAs"] != email {
		t.Errorf("Expected %v,%v,%v, but got %v,%v,%v",
			article, version, email, calledWith["article"], calledWith["version"], calledWith["loggedInAs"])
	}
}

func TestUpdateVersionNotLoggedIn(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
		return nil
	}

	// set the logged-in user
	loggedInUser = nil

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
	if w.Code != 403 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

	// check the service mock
	servMock.Mock.AssertCalled(t, "UpdateVersion", 0)
}

func TestUpdateVersionBadParam(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
		return nil
	}

	// set the logged-in user
	email := "zoo@mail.com"
	loggedInUser = &email

	// set request url
	const article int64 = 42
	const version int64 = 123456
	url := fmt.Sprintf("/articles/a/versions/%d", version)

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
	servMock.Mock.AssertCalled(t, "UpdateVersion", 0)
}

func TestUpdateVersionBadParam2(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
		return nil
	}

	// set the logged-in user
	email := "zoo@mail.com"
	loggedInUser = &email

	// set request url
	const article int64 = 42
	const version int64 = 123456
	url := fmt.Sprintf("/articles/%d/versions/a", article)

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
	servMock.Mock.AssertCalled(t, "UpdateVersion", 0)
}

func TestUpdateVersionFail(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
		return errors.New("oh no, the version couldn't be updated")
	}

	// set the logged-in user
	email := "zoo@mail.com"
	loggedInUser = &email

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
	servMock.Mock.AssertCalled(t, "UpdateVersion", 1)
	// check the service mock params manually, only checks the non-pointer parameters
	calledWith := servMock.Mock.Params["UpdateVersion"]
	if calledWith["article"] != article ||
		calledWith["version"] != version ||
		calledWith["loggedInAs"] != email {
		t.Errorf("Expected %v,%v,%v, but got %v,%v,%v",
			article, version, email, calledWith["article"], calledWith["version"], calledWith["loggedInAs"])
	}
}

func TestUpdateVersionNoFile(t *testing.T) {
	localSetup()

	// define service mock behaviour
	services.UpdateVersionMock = func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
		return nil
	}

	// set request url
	const article = "42"
	const version = "123456"
	url := fmt.Sprintf("/articles/%s/versions/%s", article, version)

	tests.TestEndpoint(t, r, "POST", url, nil, http.StatusBadRequest, nil)

	// check the service mock
	servMock.Mock.AssertCalled(t, "UpdateVersion", 0)
}

func TestGetVersion(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetVersionMock = func(article int64, version int64) (models.Version, error) {
		return exampleVersion, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/1/versions/2", nil,
		http.StatusOK, exampleVersion)

	// check the service mock
	servMock.Mock.AssertCalled(t, "GetVersion", 1)
	servMock.Mock.AssertCalled(t, "GetVersionByCommitID", 0)
	servMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"article": int64(1),
		"version": int64(2),
	})
}

func TestGetVersionBadParam(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetVersionMock = func(article int64, version int64) (models.Version, error) {
		return exampleVersion, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/a/versions/2", nil,
		http.StatusBadRequest, nil)

	// check the service mock
	servMock.Mock.AssertCalled(t, "GetVersion", 0)
}

func TestGetVersionBadParam2(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetVersionMock = func(article int64, version int64) (models.Version, error) {
		return exampleVersion, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/1/versions/a", nil,
		http.StatusBadRequest, nil)

	// check the service mock
	servMock.Mock.AssertCalled(t, "GetVersion", 0)
}

func TestGetVersionInternalError(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetVersionMock = func(article int64, version int64) (models.Version, error) {
		return models.Version{}, errors.New("fake error for testing")
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/1/versions/2", nil,
		http.StatusInternalServerError, nil)

	// check the service mock
	servMock.Mock.AssertCalled(t, "GetVersion", 1)
}

func TestGetVersionByCommit(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetVersionByCommitIDMock = func(article int64, version int64, commitID string) (models.Version, error) {
		return exampleVersion, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	url := fmt.Sprintf("/articles/1/versions/2?historyID=%s", "0123456789012345678901234567890123456789")
	tests.TestEndpoint(t, r, "GET", url, nil,
		http.StatusOK, exampleVersion)

	// check the service mock
	servMock.Mock.AssertCalled(t, "GetVersionByCommitID", 1)
}

func TestGetVersionByBadCommitParam(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.GetVersionByCommitIDMock = func(article int64, version int64, commitID string) (models.Version, error) {
		return exampleVersion, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	url := fmt.Sprintf("/articles/1/versions/2?historyID=%s", "tooshort")
	tests.TestEndpoint(t, r, "GET", url, nil,
		http.StatusBadRequest, nil)

	// check the service mock
	servMock.Mock.AssertCalled(t, "GetVersionByCommitID", 0)
}

func TestListVersion(t *testing.T) {
	localSetup()

	// define mock behaviour
	services.ListVersionsMock = func(article int64) ([]models.Version, error) {
		return exampleVersionList, nil
	}

	// this should be accessible without being logged in
	loggedInUser = nil

	// execute and test the http request
	tests.TestEndpoint(t, r, "GET", "/articles/1", nil,
		http.StatusOK, exampleVersionList)

	// check the service mock
	servMock.Mock.AssertCalled(t, "ListVersions", 1)
	servMock.Mock.AssertCalledWith(t, "ListVersions", &map[string]interface{}{
		"article": int64(1),
	})
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

var exampleVersionList = []models.Version{
	exampleVersion,
	exampleVersion,
}
var exampleVersion = models.Version{
	ArticleID:      1,
	Id:             2,
	Title:          "Version Title",
	Owners:         []string{"john@mail.com"},
	Content:        "Bla Bla Bla",
	Status:         "draft",
	LatestCommitID: "0123456789012345678901234567890123456789",
}
