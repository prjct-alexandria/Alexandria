package reviewThread

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"testing"
)

var reviewThreadRepoMock repositories.ReviewThreadRepositoryMock

var serv *services.ReviewThreadService

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
	reviewThreadRepoMock = repositories.NewReviewThreadRepositoryMock()
	servVal := services.ReviewThreadService{ReviewThreadRepository: reviewThreadRepoMock}
	serv = &servVal
}

func TestStartReviewThreadSuccess(t *testing.T) {
	// Arrange
	localSetup()

	reviewId := int64(1)
	tid := int64(-1)

	repositories.CreateReviewThreadMock = func(rid int64, tid int64) (int64, error) {
		return int64(1), nil
	}

	expected := int64(1)

	// Act
	actual, err := serv.StartReviewThread(reviewId, tid)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	reviewThreadRepoMock.Mock.AssertCalledWith(t, "CreateReviewThread", &map[string]interface{}{
		"rid": reviewId,
		"tid": tid,
	})
}
func TestStartReviewThreadDbFail(t *testing.T) {
	// Arrange
	localSetup()

	reviewId := int64(1)
	tid := int64(-1)

	repositories.CreateReviewThreadMock = func(rid int64, tid int64) (int64, error) {
		return int64(1), errors.New("error")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartReviewThread(reviewId, tid)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	reviewThreadRepoMock.Mock.AssertCalledWith(t, "CreateReviewThread", &map[string]interface{}{
		"rid": reviewId,
		"tid": tid,
	})
}
