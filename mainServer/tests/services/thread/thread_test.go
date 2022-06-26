package thread

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"testing"
)

var threadRepoMock repositories.ThreadRepositoryMock

var serv *services.ThreadService

func TestMain(m *testing.M) {
	globalSetup()
	gin.SetMode(gin.TestMode)
	m.Run()
}

// globalSetup should be called once, before any test in this file starts
func globalSetup() {

}

func localSetup() {
	// Make a clean controller with clean mocks
	threadRepoMock = repositories.NewThreadRepositoryMock()
	servVal := services.ThreadService{ThreadRepository: threadRepoMock}
	serv = &servVal
}

func TestStartThreadSuccess(t *testing.T) {
	// Arrange
	localSetup()

	threadModel := models.Thread{
		Id:         0,
		ArticleId:  3,
		SpecificId: "",
		Comments:   nil,
		Selection:  "",
	}

	articleId := int64(3)

	repositories.CreateThreadMock = func(aid int64) (int64, error) {
		return int64(1), nil
	}

	expected := int64(1)

	// Act
	actual, err := serv.StartThread(threadModel, articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)

	threadRepoMock.Mock.AssertCalledWith(t, "CreateThread", &map[string]interface{}{
		"aid": articleId,
	})
}

func TestStartThreadNoIdMatchFail(t *testing.T) {
	// Arrange
	localSetup()

	threadModel := models.Thread{
		Id:         0,
		ArticleId:  2,
		SpecificId: "",
		Comments:   nil,
		Selection:  "",
	}
	articleId := int64(3)

	expected := int64(-1)
	// Act
	actual, err := serv.StartThread(threadModel, articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, err, nil)

	threadRepoMock.Mock.AssertCalled(t, "CreateThread", 0)
}

func TestStartThreadDbFail(t *testing.T) {
	// Arrange
	localSetup()

	threadModel := models.Thread{
		Id:         0,
		ArticleId:  3,
		SpecificId: "",
		Comments:   nil,
		Selection:  "",
	}

	articleId := int64(3)

	repositories.CreateThreadMock = func(aid int64) (int64, error) {
		return int64(1), errors.New("error")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartThread(threadModel, articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, err, nil)

	threadRepoMock.Mock.AssertCalledWith(t, "CreateThread", &map[string]interface{}{
		"aid": articleId,
	})
}
