package comment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/entities"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"testing"
)

var commentRepoMock repositories.CommentRepositoryMock

var serv *services.CommentService

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
	commentRepoMock = repositories.NewCommentRepositoryMock()
	servVal := services.CommentService{CommentRepository: commentRepoMock}
	serv = &servVal
}

func TestSaveCommentSuccess(t *testing.T) {
	localSetup()

	// Arrange
	threadId := int64(3)

	comment := entities.Comment{
		Id:           0,
		AuthorId:     "johndoe@gmail.com",
		ThreadId:     threadId,
		Content:      "I like this app",
		CreationDate: "1656200000",
	}

	loggedInAs := "johndoe@gmail.com"

	repositories.SaveCommentMock = func(id entities.Comment) (int64, error) {
		return int64(1), nil
	}

	expected := int64(1)

	// Act
	actual, err := serv.SaveComment(comment, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	commentRepoMock.Mock.AssertCalledWith(t, "SaveComment", &map[string]interface{}{
		"id": comment,
	})
}

func TestSaveCommentAuthFail(t *testing.T) {
	localSetup()

	// Arrange
	threadId := int64(3)

	comment := entities.Comment{
		Id:           0,
		AuthorId:     "johndoe@gmail.com",
		ThreadId:     threadId,
		Content:      "I like this app",
		CreationDate: "1656200000",
	}

	loggedInAs := "impostor@gmail.com"

	repositories.SaveCommentMock = func(id entities.Comment) (int64, error) {
		return int64(1), nil
	}

	expected := int64(-1)

	// Act
	actual, err := serv.SaveComment(comment, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	commentRepoMock.Mock.AssertCalled(t, "SaveComment", 0)
}

func TestSaveCommentDbFail(t *testing.T) {
	localSetup()

	// Arrange
	threadId := int64(3)

	comment := entities.Comment{
		Id:           0,
		AuthorId:     "johndoe@gmail.com",
		ThreadId:     threadId,
		Content:      "I like this app",
		CreationDate: "1656200000",
	}

	loggedInAs := "johndoe@gmail.com"

	repositories.SaveCommentMock = func(id entities.Comment) (int64, error) {
		return int64(1), errors.New("there is one impostor among us")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.SaveComment(comment, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	commentRepoMock.Mock.AssertCalledWith(t, "SaveComment", &map[string]interface{}{
		"id": comment,
	})
}
