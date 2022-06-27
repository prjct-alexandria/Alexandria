package requestThread

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	"strconv"
	"testing"
)

var requestThreadRepoMock repositories.RequestThreadRepositoryMock
var versionRepoMock repositories.VersionRepositoryMock
var requestRepoMock repositories.RequestRepositoryMock

var serv *services.RequestThreadService

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
	requestThreadRepoMock = repositories.NewRequestThreadRepositoryMock()
	versionRepoMock = repositories.NewVersionRepositoryMock()
	requestRepoMock = repositories.NewRequestRepositoryMock()
	servVal := services.RequestThreadService{
		RequestThreadRepository: requestThreadRepoMock,
		VersionRepository:       versionRepoMock,
		RequestRepository:       requestRepoMock,
	}
	serv = &servVal
}

func TestStartRequestThreadSuccess(t *testing.T) {
	// Assert
	localSetup()

	email1 := "JaneDoe@gmail.com"
	email2 := "JohnDoe@gmail.com"
	requestId := int64(1)
	threadId := int64(0)

	loggedInAs := email1
	sourceId := int64(1)
	targetId := int64(2)

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return entities.Request{
			RequestID:       request,
			ArticleID:       int64(1),
			SourceVersionID: sourceId,
			SourceHistoryID: "",
			TargetVersionID: targetId,
			TargetHistoryID: "",
			Status:          "pending",
			Conflicted:      false,
		}, nil
	}

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceId && email == email1 {
			return true, nil
		} else if version == targetId && email == email2 {
			return true, nil
		} else {
			return false, nil
		}
	}

	repositories.CreateRequestThreadMock = func(rid int64, tid int64) (int64, error) {
		return int64(5), nil
	}

	expected := int64(5)

	// Act
	actual, err := serv.StartRequestThread(requestId, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetId,
		"email":   loggedInAs,
	})

	requestThreadRepoMock.Mock.AssertCalledWith(t, "CreateRequestThread", &map[string]interface{}{
		"rid": requestId,
		"tid": threadId,
	})
}

func TestStartRequestThreadDbReqFail(t *testing.T) {
	// Assert
	localSetup()

	email1 := "JaneDoe@gmail.com"
	requestId := int64(1)
	threadId := int64(0)
	loggedInAs := email1

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return entities.Request{}, errors.New("error")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartRequestThread(requestId, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 0)

	requestThreadRepoMock.Mock.AssertCalled(t, "CreateRequestThread", 0)
}

func TestStartRequestThreadDbSourceFail(t *testing.T) {
	// Assert
	localSetup()

	email1 := "JaneDoe@gmail.com"
	email2 := "JohnDoe@gmail.com"
	requestId := int64(1)
	threadId := int64(0)

	loggedInAs := email1
	sourceId := int64(1)
	targetId := int64(2)

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return entities.Request{
			RequestID:       request,
			ArticleID:       int64(1),
			SourceVersionID: sourceId,
			SourceHistoryID: "",
			TargetVersionID: targetId,
			TargetHistoryID: "",
			Status:          "pending",
			Conflicted:      false,
		}, nil
	}

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceId && email == email1 {
			return true, errors.New("error")
		} else if version == targetId && email == email2 {
			return true, nil
		} else {
			return false, nil
		}
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartRequestThread(requestId, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 1)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": sourceId,
		"email":   loggedInAs,
	})

	requestThreadRepoMock.Mock.AssertCalled(t, "CreateRequestThread", 0)
}

func TestStartRequestThreadDbTargetFail(t *testing.T) {
	// Assert
	localSetup()

	email1 := "JaneDoe@gmail.com"
	email2 := "JohnDoe@gmail.com"
	requestId := int64(1)
	threadId := int64(0)

	loggedInAs := email1
	sourceId := int64(1)
	targetId := int64(2)

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return entities.Request{
			RequestID:       request,
			ArticleID:       int64(1),
			SourceVersionID: sourceId,
			SourceHistoryID: "",
			TargetVersionID: targetId,
			TargetHistoryID: "",
			Status:          "pending",
			Conflicted:      false,
		}, nil
	}

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceId && email == email1 {
			return true, nil
		} else if version == targetId && email == email2 {
			return true, errors.New("error")
		} else {
			return false, errors.New("error")
		}
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartRequestThread(requestId, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetId,
		"email":   loggedInAs,
	})

	requestThreadRepoMock.Mock.AssertCalled(t, "CreateRequestThread", 0)
}

func TestStartRequestThreadNotAllowedFail(t *testing.T) {
	// Assert
	localSetup()

	email1 := "JaneDoe@gmail.com"
	email2 := "JohnDoe@gmail.com"
	requestId := int64(1)
	threadId := int64(0)

	loggedInAs := "impostor@gmail.com"
	sourceId := int64(1)
	targetId := int64(2)

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return entities.Request{
			RequestID:       request,
			ArticleID:       int64(1),
			SourceVersionID: sourceId,
			SourceHistoryID: "",
			TargetVersionID: targetId,
			TargetHistoryID: "",
			Status:          "pending",
			Conflicted:      false,
		}, nil
	}

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceId && email == email1 {
			return true, nil
		} else if version == targetId && email == email2 {
			return true, nil
		} else {
			return false, nil
		}
	}

	repositories.CreateRequestThreadMock = func(rid int64, tid int64) (int64, error) {
		return int64(5), nil
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartRequestThread(requestId, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetId,
		"email":   loggedInAs,
	})

	requestThreadRepoMock.Mock.AssertCalled(t, "CreateRequestThread", 0)
}

func TestStartRequestThreadDbFail(t *testing.T) {
	// Assert
	localSetup()

	email1 := "JaneDoe@gmail.com"
	email2 := "JohnDoe@gmail.com"
	requestId := int64(1)
	threadId := int64(0)

	loggedInAs := email1
	sourceId := int64(1)
	targetId := int64(2)

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return entities.Request{
			RequestID:       request,
			ArticleID:       int64(1),
			SourceVersionID: sourceId,
			SourceHistoryID: "",
			TargetVersionID: targetId,
			TargetHistoryID: "",
			Status:          "pending",
			Conflicted:      false,
		}, nil
	}

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceId && email == email1 {
			return true, nil
		} else if version == targetId && email == email2 {
			return true, nil
		} else {
			return false, nil
		}
	}

	repositories.CreateRequestThreadMock = func(rid int64, tid int64) (int64, error) {
		return int64(5), errors.New("error")
	}

	expected := int64(-1)

	// Act
	actual, err := serv.StartRequestThread(requestId, threadId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetId,
		"email":   loggedInAs,
	})

	requestThreadRepoMock.Mock.AssertCalledWith(t, "CreateRequestThread", &map[string]interface{}{
		"rid": requestId,
		"tid": threadId,
	})
}

func TestGetRequestThreadsSuccess(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	requestId := int64(2)

	threadList := []models.Thread{
		{
			Id:         0,
			ArticleId:  articleId,
			SpecificId: strconv.FormatInt(requestId, 10),
			Comments:   nil,
		}, {
			Id:         1,
			ArticleId:  articleId,
			SpecificId: strconv.FormatInt(requestId, 10),
			Comments:   nil,
		},
	}

	repositories.GetRequestThreadsMock = func(aid int64, rid int64) ([]models.Thread, error) {
		return threadList, nil
	}

	expected := threadList

	// Act
	actual, err := serv.GetRequestThreads(articleId, requestId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	requestThreadRepoMock.Mock.AssertCalledWith(t, "GetRequestThreads", &map[string]interface{}{
		"aid": articleId,
		"rid": requestId,
	})
}

func TestGetRequestThreadsFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	requestId := int64(2)

	threadList := []models.Thread{
		{
			Id:         0,
			ArticleId:  articleId,
			SpecificId: strconv.FormatInt(requestId, 10),
			Comments:   nil,
		}, {
			Id:         1,
			ArticleId:  articleId,
			SpecificId: strconv.FormatInt(requestId, 10),
			Comments:   nil,
		},
	}

	repositories.GetRequestThreadsMock = func(aid int64, rid int64) ([]models.Thread, error) {
		return threadList, errors.New("error")
	}

	var expected []models.Thread

	// Act
	actual, err := serv.GetRequestThreads(articleId, requestId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestThreadRepoMock.Mock.AssertCalledWith(t, "GetRequestThreads", &map[string]interface{}{
		"aid": articleId,
		"rid": requestId,
	})
}
