package commitSelectionThread

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"testing"
)

var commitSelectionThreadRepoMock repositories.CommitSelectionThreadRepositoryMock

var serv *services.CommitSelectionThreadService

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
	// Make a clean service with clean mocks
	commitSelectionThreadRepoMock = repositories.NewCommitSelectionThreadRepositoryMock()
	servVal := services.CommitSelectionThreadService{CommitSelectionThreadRepository: commitSelectionThreadRepoMock}
	serv = &servVal
}

func TestStartCommitSelectionThreadSuccess(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	tid := int64(-1)
	section := "some article section: true"

	repositories.CreateCommitSelectionThreadMock = func(cid string, tid int64, section string) (int64, error) {
		return int64(1), nil
	}

	expected := int64(1)

	// Act
	actual, err := serv.StartCommitSelectionThread(commitId, tid, section)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	commitSelectionThreadRepoMock.Mock.AssertCalledWith(t, "CreateCommitSelectionThread", &map[string]interface{}{
		"cid":     commitId,
		"tid":     tid,
		"section": section,
	})
}

func TestStartCommitSelectionThreadNoCommitIdFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "abcdefghijklmnop12345"
	tid := int64(-1)
	section := "some article section: true"

	expected := int64(-1)

	// Act
	actual, err := serv.StartCommitSelectionThread(commitId, tid, section)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	commitSelectionThreadRepoMock.Mock.AssertCalled(t, "CreateCommitSelectionThread", 0)
}

func TestStartCommitSelectionThreadDbFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	tid := int64(-1)
	section := "some article section: true"

	repositories.CreateCommitSelectionThreadMock = func(cid string, tid int64, section string) (int64, error) {
		return int64(1), errors.New("oops an error has occurred, sorry not sorry")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartCommitSelectionThread(commitId, tid, section)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	commitSelectionThreadRepoMock.Mock.AssertCalledWith(t, "CreateCommitSelectionThread", &map[string]interface{}{
		"cid":     commitId,
		"tid":     tid,
		"section": section,
	})
}

func TestGetCommitSelectionThreadsSuccess(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	articleId := int64(1)

	threadList := []models.SelectionThread{
		{
			Id:         0,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
			Selection:  "",
		}, {
			Id:         1,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
			Selection:  "",
		},
	}

	repositories.GetCommitSelectionThreadsMock = func(aid int64, cid string) ([]models.SelectionThread, error) {
		return threadList, nil
	}

	expected := threadList

	// Act
	actual, err := serv.GetCommitSelectionThreads(commitId, articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	commitSelectionThreadRepoMock.Mock.AssertCalledWith(t, "GetCommitSelectionThreads", &map[string]interface{}{
		"aid": articleId,
		"cid": commitId,
	})
}

func TestGetCommitSelectionThreadsCommitIdFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "abcd217361782637812zyxwvu"
	articleId := int64(1)

	var expected []models.SelectionThread

	// Act
	actual, err := serv.GetCommitSelectionThreads(commitId, articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	commitSelectionThreadRepoMock.Mock.AssertCalled(t, "GetCommitSelectionThreads", 0)
}

func TestGetCommitSelectionThreadsDbFail(t *testing.T) {
	// Arrange
	localSetup()

	commitId := "002f3e351edfc91bb934377b007bf404f3bdb78a"
	articleId := int64(1)

	threadList := []models.SelectionThread{
		{
			Id:         0,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
			Selection:  "",
		}, {
			Id:         1,
			ArticleId:  articleId,
			SpecificId: commitId,
			Comments:   nil,
			Selection:  "",
		},
	}

	repositories.GetCommitSelectionThreadsMock = func(aid int64, cid string) ([]models.SelectionThread, error) {
		return threadList, errors.New("I am an error, mwahahaha")
	}

	var expected []models.SelectionThread

	// Act
	actual, err := serv.GetCommitSelectionThreads(commitId, articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	commitSelectionThreadRepoMock.Mock.AssertCalledWith(t, "GetCommitSelectionThreads", &map[string]interface{}{
		"aid": articleId,
		"cid": commitId,
	})
}
