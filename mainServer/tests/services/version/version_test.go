package version

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	mocks "mainServer/tests/util"
	"testing"
)

var versionRepoMock repositories.VersionRepositoryMock
var storerRepoMock repositories.StorerMock
var userRepoMock repositories.UserRepositoryMock

var serv *services.VersionService

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
	userRepoMock = repositories.NewUserRepositoryMock()
	storerRepoMock = repositories.StorerMock{Mock: mocks.NewMock()}
	versionRepoMock = repositories.NewVersionRepositoryMock()
	servVal := services.VersionService{
		VersionRepo: versionRepoMock,
		Storer:      storerRepoMock,
		UserRepo:    userRepoMock,
	}
	serv = &servVal
}

//TODO: Add bad weather cases for all tests

func TestListVersionSuccess(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(2)

	versionEntity := entities.Version{
		ArticleID:      articleId,
		Id:             int64(2),
		Title:          "Lorem Ipsum",
		Owners:         nil,
		Status:         "draft",
		LatestCommitID: "",
	}

	repositories.GetVersionsByArticleMock = func(article int64) ([]entities.Version, error) {
		return []entities.Version{versionEntity}, nil
	}

	expected := []models.Version{
		{
			ArticleID:      versionEntity.ArticleID,
			Id:             versionEntity.Id,
			Title:          versionEntity.Title,
			Owners:         versionEntity.Owners,
			Content:        "",
			Status:         versionEntity.Status,
			LatestCommitID: versionEntity.LatestCommitID,
		},
	}

	// Act
	actual, err := serv.ListVersions(articleId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	versionRepoMock.Mock.AssertCalledWith(t, "GetVersionsByArticle", &map[string]interface{}{
		"article": articleId,
	})
}

func TestGetVersionSuccess(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	versionId := int64(2)

	versionEntity := entities.Version{
		ArticleID:      articleId,
		Id:             versionId,
		Title:          "Lorem Ipsum",
		Owners:         nil,
		Status:         "draft",
		LatestCommitID: "",
	}

	repositories.GetVersionMockStorer = func(article int64, version int64) (string, error) {
		if version == versionEntity.Id {
			return "contentHello", nil
		}
		return "", nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		return versionEntity, nil
	}

	expected := models.Version{
		ArticleID:      versionEntity.ArticleID,
		Id:             versionEntity.Id,
		Title:          versionEntity.Title,
		Owners:         versionEntity.Owners,
		Content:        "contentHello",
		Status:         versionEntity.Status,
		LatestCommitID: versionEntity.LatestCommitID,
	}

	// Act
	actual, err := serv.GetVersion(articleId, versionId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	storerRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"article": articleId,
		"version": versionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": versionId,
	})

}

func TestGetVersionByCommitIdSuccess(t *testing.T) {
	// Arrange

	localSetup()

	articleId := int64(1)
	versionId := int64(2)
	commitId := "helloWorldCommitId"

	versionEntity := entities.Version{
		ArticleID:      articleId,
		Id:             versionId,
		Title:          "Lorem Ipsum",
		Owners:         nil,
		Status:         "draft",
		LatestCommitID: "",
	}

	repositories.GetVersionByCommitMock = func(article int64, commit string) (string, error) {
		return "contentHello", nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		return versionEntity, nil
	}

	expected := models.Version{
		ArticleID:      versionEntity.ArticleID,
		Id:             versionEntity.Id,
		Title:          versionEntity.Title,
		Owners:         versionEntity.Owners,
		Content:        "contentHello",
		Status:         versionEntity.Status,
		LatestCommitID: versionEntity.LatestCommitID,
	}

	// Act
	actual, err := serv.GetVersionByCommitID(articleId, versionId, commitId)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	storerRepoMock.Mock.AssertCalledWith(t, "GetVersionByCommit", &map[string]interface{}{
		"article": articleId,
		"commit":  commitId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": versionId,
	})
}

func TestCreateVersionFromSuccess(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	versionId := int64(2)
	commitId := "thisIsACommitId"

	newTitle := "newTitle"
	newOwnersDupl := []string{"johndoe@gmail.com", "janedoe@gmail.com", "johndoe@gmail.com"}
	newOwners := []string{"johndoe@gmail.com", "janedoe@gmail.com"}
	loggedInAs := "johndoe@gmail.com"
	newVersionId := int64(4)

	versionEntity := entities.Version{
		ArticleID:      articleId,
		Id:             versionId,
		Title:          "Lorem Ipsum",
		Owners:         nil,
		Status:         "draft",
		LatestCommitID: "",
	}

	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateVersionMock = func(version entities.Version) (entities.Version, error) {
		versionEntity.Id = newVersionId
		versionEntity.Title = newTitle
		versionEntity.Owners = newOwners
		return versionEntity, nil
	}

	repositories.CreateVersionFromMock = func(article int64, source int64, target int64) (string, error) {
		return commitId, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return nil
	}

	expected := models.Version{
		ArticleID:      articleId,
		Id:             newVersionId,
		Title:          newTitle,
		Owners:         newOwners,
		Content:        "",
		LatestCommitID: commitId,
	}

	// Act
	actual, err := serv.CreateVersionFrom(articleId, versionId, newTitle, newOwnersDupl, loggedInAs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": newOwners[1],
	})

	versionRepoMock.Mock.AssertCalledWith(t, "CreateVersion", &map[string]interface{}{
		"version": entities.Version{
			ArticleID: articleId,
			Title:     newTitle,
			Owners:    newOwners,
		},
	})

	storerRepoMock.Mock.AssertCalledWith(t, "CreateVersionFrom", &map[string]interface{}{
		"article": articleId,
		"source":  versionId,
		"target":  newVersionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "UpdateVersionLatestCommit", &map[string]interface{}{
		"version": newVersionId,
		"commit":  commitId,
	})
}

func TestUpdateVersion(t *testing.T) {
	// TODO
}

func TestGetVersionFiles(t *testing.T) {
	// TODO
}
