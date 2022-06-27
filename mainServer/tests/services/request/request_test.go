package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	mocks "mainServer/tests/util"
	"testing"
)

var requestRepoMock repositories.RequestRepositoryMock
var versionRepoMock repositories.VersionRepositoryMock
var storerMock repositories.StorerMock

var serv *services.RequestService

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
	requestRepoMock = repositories.NewRequestRepositoryMock()
	versionRepoMock = repositories.NewVersionRepositoryMock()
	storerMock = repositories.StorerMock{Mock: mocks.NewMock()}
	servVal := services.RequestService{Repo: requestRepoMock, Versionrepo: versionRepoMock, Storer: storerMock}
	serv = &servVal
}

// Test create request

func TestCreateRequestSuccessSourceOwner(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	sourceVersionId := int64(3)
	targetVersionId := int64(4)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	loggedInAs := "johnDoe@gmail.com"

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceVersionId && email == loggedInAs {
			return true, nil
		}
		return false, nil
	}

	repositories.CreateRequestMock = func(req entities.Request) (entities.Request, error) {
		return req, nil
	}

	expected := models.Request(req)

	// Act
	actual, err := serv.CreateRequest(articleId, sourceVersionId, targetVersionId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetVersionId,
		"email":   loggedInAs,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "CreateRequest", &map[string]interface{}{
		"req": req,
	})
}

func TestCreateRequestSuccessTargetOwner(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	sourceVersionId := int64(3)
	targetVersionId := int64(4)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	loggedInAs := "johnDoe@gmail.com"

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == targetVersionId && email == loggedInAs {
			return true, nil
		}
		return false, nil
	}

	repositories.CreateRequestMock = func(req entities.Request) (entities.Request, error) {
		return req, nil
	}

	expected := models.Request(req)

	// Act
	actual, err := serv.CreateRequest(articleId, sourceVersionId, targetVersionId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetVersionId,
		"email":   loggedInAs,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "CreateRequest", &map[string]interface{}{
		"req": req,
	})
}

func TestCreateRequestDbSourceOwnerCheckFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	sourceVersionId := int64(3)
	targetVersionId := int64(4)

	loggedInAs := "johnDoe@gmail.com"

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceVersionId && email == loggedInAs {
			return true, errors.New("error")
		}
		return false, nil
	}

	expected := models.Request{}

	// Act
	actual, err := serv.CreateRequest(articleId, sourceVersionId, targetVersionId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 1)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": sourceVersionId,
		"email":   loggedInAs,
	})
}

func TestCreateRequestDbTargetOwnerCheckFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	sourceVersionId := int64(3)
	targetVersionId := int64(4)

	loggedInAs := "johnDoe@gmail.com"

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == targetVersionId && email == loggedInAs {
			return true, errors.New("error")
		}
		return false, nil
	}

	expected := models.Request{}

	// Act
	actual, err := serv.CreateRequest(articleId, sourceVersionId, targetVersionId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetVersionId,
		"email":   loggedInAs,
	})
}

func TestCreateRequestNotAnOwnerFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	sourceVersionId := int64(3)
	targetVersionId := int64(4)

	loggedInAs := "johnDoe@gmail.com"

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		return false, nil
	}

	expected := models.Request{}

	// Act
	actual, err := serv.CreateRequest(articleId, sourceVersionId, targetVersionId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetVersionId,
		"email":   loggedInAs,
	})
}

func TestCreateRequestDbFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	sourceVersionId := int64(3)
	targetVersionId := int64(4)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	loggedInAs := "johnDoe@gmail.com"

	repositories.CheckIfOwnerMock = func(version int64, email string) (bool, error) {
		if version == sourceVersionId && email == loggedInAs {
			return true, nil
		}
		return false, nil
	}

	repositories.CreateRequestMock = func(req entities.Request) (entities.Request, error) {
		return req, errors.New("error")
	}

	expected := models.Request{}

	// Act
	actual, err := serv.CreateRequest(articleId, sourceVersionId, targetVersionId, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	versionRepoMock.Mock.AssertCalled(t, "CheckIfOwner", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "CheckIfOwner", &map[string]interface{}{
		"version": targetVersionId,
		"email":   loggedInAs,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "CreateRequest", &map[string]interface{}{
		"req": req,
	})
}

// Test reject request

func TestRejectRequestSuccess(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	repositories.SetStatusMock = func(request int64, status string) error {
		return nil
	}

	// Act
	actual := serv.RejectRequest(requestId, loggedInAs)

	// Assert
	assert.Equal(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "SetStatus", &map[string]interface{}{
		"request": requestId,
		"status":  entities.RequestRejected,
	})
}

func TestRejectRequestDbFetchFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, errors.New("error")
	}

	// Act
	actual := serv.RejectRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 0)

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)
}

func TestRejectRequestSourceVersionFetchFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, errors.New("error")
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	// Act
	actual := serv.RejectRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 1)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": sourceVersionId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)
}

func TestRejectRequestTargetVersionFetchFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, errors.New("error")
		}
		return entities.Version{}, nil
	}

	// Act
	actual := serv.RejectRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)
}

func TestRejectRequestNotAnOwnerFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "impostor@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", "johnDoe@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	// Act
	actual := serv.RejectRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)
}

func TestRejectRequestSetStatusFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	repositories.SetStatusMock = func(request int64, status string) error {
		return errors.New("hello")
	}

	// Act
	actual := serv.RejectRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "SetStatus", &map[string]interface{}{
		"request": requestId,
		"status":  entities.RequestRejected,
	})
}

// Test accept request

func TestAcceptRequestSuccess(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	commitId := "thisIsACommitId"

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	repositories.MergeMock = func(article int64, source int64, target int64) (string, error) {
		return commitId, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return nil
	}

	repositories.SetStatusMock = func(request int64, status string) error {
		return nil
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.Equal(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "SetStatus", &map[string]interface{}{
		"request": requestId,
		"status":  entities.RequestAccepted,
	})

	storerMock.Mock.AssertCalledWith(t, "Merge", &map[string]interface{}{
		"article": articleId,
		"source":  sourceVersionId,
		"target":  targetVersionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "UpdateVersionLatestCommit", &map[string]interface{}{
		"version": targetVersionId,
		"commit":  commitId,
	})
}

func TestAcceptRequestDbFetchFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, errors.New("error")
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 0)

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)

	storerMock.Mock.AssertCalled(t, "Merge", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestAcceptRequestSourceVersionFetchFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, errors.New("error")
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 1)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": sourceVersionId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)

	storerMock.Mock.AssertCalled(t, "Merge", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestAcceptRequestTargetVersionFetchFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, errors.New("error")
		}
		return entities.Version{}, nil
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)

	storerMock.Mock.AssertCalled(t, "Merge", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestAcceptRequestNotAnOwnerFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "impostor@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", "johnDoe@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)

	storerMock.Mock.AssertCalled(t, "Merge", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestAcceptRequestSetStatusFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)
	commitId := "thisIsACommitId"

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return nil
	}

	repositories.MergeMock = func(article int64, source int64, target int64) (string, error) {
		return commitId, nil
	}

	repositories.SetStatusMock = func(request int64, status string) error {
		return errors.New("hello")
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	storerMock.Mock.AssertCalledWith(t, "Merge", &map[string]interface{}{
		"article": articleId,
		"source":  sourceVersionId,
		"target":  targetVersionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "UpdateVersionLatestCommit", &map[string]interface{}{
		"version": targetVersionId,
		"commit":  commitId,
	})

	requestRepoMock.Mock.AssertCalledWith(t, "SetStatus", &map[string]interface{}{
		"request": requestId,
		"status":  entities.RequestAccepted,
	})
}

func TestAcceptRequestMergeFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)
	commitId := "thisIsACommitId"

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return nil
	}

	repositories.MergeMock = func(article int64, source int64, target int64) (string, error) {
		return commitId, errors.New("thisIsAnAerror")
	}

	repositories.SetStatusMock = func(request int64, status string) error {
		return nil
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	storerMock.Mock.AssertCalledWith(t, "Merge", &map[string]interface{}{
		"article": articleId,
		"source":  sourceVersionId,
		"target":  targetVersionId,
	})

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)
}

func TestAcceptRequestUpdateCommitFail(t *testing.T) {
	// Arrange
	localSetup()

	requestId := int64(42)
	loggedInAs := "johnDoe@gmail.com"
	articleId := int64(1)
	sourceVersionId := int64(5)
	targetVersionId := int64(12)
	commitId := "thisIsACommitId"

	req := entities.Request{
		ArticleID:       articleId,
		SourceVersionID: sourceVersionId,
		TargetVersionID: targetVersionId,
	}

	repositories.GetRequestMock = func(request int64) (entities.Request, error) {
		return req, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == sourceVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             sourceVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com"},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		} else if version == targetVersionId {
			return entities.Version{
				ArticleID:      articleId,
				Id:             targetVersionId,
				Title:          "title",
				Owners:         []string{"jane@gmail.com", loggedInAs},
				Status:         "pending",
				LatestCommitID: "",
			}, nil
		}
		return entities.Version{}, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return errors.New("hello")
	}

	repositories.MergeMock = func(article int64, source int64, target int64) (string, error) {
		return commitId, nil
	}

	repositories.SetStatusMock = func(request int64, status string) error {
		return nil
	}

	// Act
	actual := serv.AcceptRequest(requestId, loggedInAs)

	// Assert
	assert.NotEqual(t, nil, actual)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequest", &map[string]interface{}{
		"request": requestId,
	})

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": targetVersionId,
	})

	storerMock.Mock.AssertCalledWith(t, "Merge", &map[string]interface{}{
		"article": articleId,
		"source":  sourceVersionId,
		"target":  targetVersionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "UpdateVersionLatestCommit", &map[string]interface{}{
		"version": targetVersionId,
		"commit":  commitId,
	})

	requestRepoMock.Mock.AssertCalled(t, "SetStatus", 0)
}

// Test GetRequest (with comparison)

func TestGetRequestSuccess(t *testing.T) {

}

//Test UpdateRequestComparison

func TestUpdateRequestComparisonSuccess(t *testing.T) {

}

// Test GetRequestList

func TestGetRequestListSuccess(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(2)
	sourceId := int64(-1)
	targetId := int64(5)
	relatedId := int64(-1)

	expected := []models.RequestListElement{
		{
			Request: models.Request{
				RequestID:       7,
				ArticleID:       articleId,
				SourceVersionID: 4,
				SourceHistoryID: "",
				TargetVersionID: targetId,
				TargetHistoryID: "",
				Status:          "pending",
				Conflicted:      false,
			},
			SourceTitle: "abc",
			TargetTitle: "abc"},
	}

	repositories.GetRequestListMock = func(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error) {
		return expected, nil
	}

	// Act
	actual, err := serv.GetRequestList(articleId, sourceId, targetId, relatedId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequestList", &map[string]interface{}{
		"articleId": articleId,
		"sourceId":  sourceId,
		"targetId":  targetId,
		"relatedId": relatedId,
	})
}

func TestGetRequestListFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(2)
	sourceId := int64(-1)
	targetId := int64(5)
	relatedId := int64(-1)

	expected := []models.RequestListElement{
		{
			Request: models.Request{
				RequestID:       7,
				ArticleID:       articleId,
				SourceVersionID: 4,
				SourceHistoryID: "",
				TargetVersionID: targetId,
				TargetHistoryID: "",
				Status:          "pending",
				Conflicted:      false,
			},
			SourceTitle: "abc",
			TargetTitle: "abc"},
	}

	repositories.GetRequestListMock = func(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error) {
		return expected, errors.New("error")
	}

	// Act
	actual, err := serv.GetRequestList(articleId, sourceId, targetId, relatedId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, nil, err)

	requestRepoMock.Mock.AssertCalledWith(t, "GetRequestList", &map[string]interface{}{
		"articleId": articleId,
		"sourceId":  sourceId,
		"targetId":  targetId,
		"relatedId": relatedId,
	})
}
