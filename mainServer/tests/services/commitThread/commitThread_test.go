package commitThread

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"testing"
)

var commitThreadRepoMock repositories.CommitThreadRepositoryMock

var serv *services.CommitThreadService

func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test in this file starts
func globalSetup() {

}

// localSetup should be called before each individual test
func localSetup() {
	// Make a clean controller with clean mocks
	commitThreadRepoMock = repositories.NewCommitThreadRepositoryMock()
	servVal := services.CommitThreadService{CommitThreadRepository: commitThreadRepoMock}
	serv = &servVal
}

func TestStartCommitThreadSuccess(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	tid := int64(-1)

	repositories.CreateCommitThreadMock = func(cid string, tid int64) (int64, error) {
		return int64(1), nil
	}

	expected := int64(1)

	// Act
	actual, err := serv.StartCommitThread(commitId, tid)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)

	commitThreadRepoMock.Mock.AssertCalledWith(t, "CreateCommitThread", &map[string]interface{}{
		"cid": commitId,
		"tid": tid,
	})
}

func TestStartCommitThreadNoCommitIdFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "abcdefghijklmnop12345"
	tid := int64(-1)

	expected := int64(-1)

	// Act
	actual, err := serv.StartCommitThread(commitId, tid)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, err, nil)

	commitThreadRepoMock.Mock.AssertCalled(t, "CreateCommitThread", 0)
}

func TestStartCommitThreadDbFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	tid := int64(-1)

	repositories.CreateCommitThreadMock = func(cid string, tid int64) (int64, error) {
		return int64(1), errors.New("oops an error has occurred, sorry not sorry")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartCommitThread(commitId, tid)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, err, nil)

	commitThreadRepoMock.Mock.AssertCalledWith(t, "CreateCommitThread", &map[string]interface{}{
		"cid": commitId,
		"tid": tid,
	})
}

func TestGetCommitThreadsSuccess(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	articleId := int64(1)

	threadList := []models.Thread{
		{
			Id:         0,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
		}, {
			Id:         1,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
		},
	}

	repositories.GetCommitThreadsMock = func(aid int64, cid string) ([]models.Thread, error) {
		return threadList, nil
	}

	expected := threadList

	// Act
	actual, err := serv.GetCommitThreads(articleId, commitId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)

	commitThreadRepoMock.Mock.AssertCalledWith(t, "GetCommitThreads", &map[string]interface{}{
		"aid": articleId,
		"cid": commitId,
	})
}

func TestGetCommitThreadsCommitIdFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "abcd217361782637812zyxwvu"
	articleId := int64(1)

	var expected []models.Thread

	// Act
	actual, err := serv.GetCommitThreads(articleId, commitId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, err, nil)

	commitThreadRepoMock.Mock.AssertCalled(t, "GetCommitThreads", 0)
}

func TestGetCommitThreadsDbFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	articleId := int64(1)

	threadList := []models.Thread{
		{
			Id:         0,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
		}, {
			Id:         1,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
		},
	}

	repositories.GetCommitThreadsMock = func(aid int64, cid string) ([]models.Thread, error) {
		return threadList, errors.New("I am an error, mwahahaha")
	}

	var expected []models.Thread

	// Act
	actual, err := serv.GetCommitThreads(articleId, commitId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, err, nil)

	commitThreadRepoMock.Mock.AssertCalledWith(t, "GetCommitThreads", &map[string]interface{}{
		"aid": articleId,
		"cid": commitId,
	})
}
