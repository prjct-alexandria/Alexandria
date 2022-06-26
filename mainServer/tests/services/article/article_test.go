package article

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/services"
	"mainServer/tests/mocks/repositories"
	mocks "mainServer/tests/util"
	"testing"
)

var articleRepoMock repositories.ArticleRepositoryMock
var versionRepoMock repositories.VersionRepositoryMock
var userRepoMock repositories.UserRepositoryMock
var storerMock repositories.StorerMock

var serv *services.ArticleService

// TestMain is a keyword function, this is run by the testing package before other tests
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
	articleRepoMock = repositories.NewArticleRepositoryMock()
	versionRepoMock = repositories.NewVersionRepositoryMock()
	userRepoMock = repositories.NewUserRepositoryMock()
	storerMock = repositories.StorerMock{Mock: mocks.NewMock()}
	servVal := services.NewArticleService(articleRepoMock, versionRepoMock, userRepoMock, storerMock)
	serv = &servVal
}

func TestGetMainVersionSucces(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)
	mainVersionId := int64(2)

	repositories.GetMainVersionMock = func(article int64) (int64, error) {
		return mainVersionId, nil
	}

	expected := mainVersionId

	// Act
	actual, err := serv.GetMainVersion(articleId)

	// Assert
	assert.Equal(t, actual, expected)
	assert.Equal(t, err, nil)

	articleRepoMock.Mock.AssertCalledWith(t, "GetMainVersion", &map[string]interface{}{
		"article": articleId,
	})
}

func TestGetMainVersionFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId := int64(1)

	repositories.GetMainVersionMock = func(article int64) (int64, error) {
		return -1, errors.New("haha I am an error")
	}

	expected := int64(-1)
	// Act
	actual, err := serv.GetMainVersion(articleId)
	// Assert
	assert.Equal(t, actual, expected)
	assert.NotEqual(t, err, nil)

	articleRepoMock.Mock.AssertCalledWith(t, "GetMainVersion", &map[string]interface{}{
		"article": articleId,
	})
}

func TestGetArticleListSuccess(t *testing.T) {
	// Arrange
	localSetup()

	articleId1 := int64(1)
	articleId2 := int64(2)
	mainVersionId1 := int64(3)
	mainVersionId2 := int64(4)
	title1 := "Lorem ipsum"
	title2 := "Dolor sit amet"
	owners1 := []string{"johndoe@gmail.com"}
	owners2 := []string{"janedoe@gmail.com"}

	repositories.GetAllArticlesMock = func() ([]entities.Article, error) {
		return []entities.Article{
			{
				Id:            articleId1,
				MainVersionID: mainVersionId1,
			}, {
				Id:            articleId2,
				MainVersionID: mainVersionId2,
			},
		}, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == mainVersionId1 {
			return entities.Version{
				ArticleID:      articleId1,
				Id:             mainVersionId1,
				Title:          title1,
				Owners:         owners1,
				Status:         "draft",
				LatestCommitID: "commitId",
			}, nil
		} else {
			return entities.Version{
				ArticleID:      articleId2,
				Id:             mainVersionId2,
				Title:          title2,
				Owners:         owners2,
				Status:         "draft",
				LatestCommitID: "commitId",
			}, nil
		}
	}

	expected := []models.ArticleListElement{{
		Id:            articleId1,
		MainVersionId: mainVersionId1,
		Title:         title1,
		Owners:        owners1},
		{
			Id:            articleId2,
			MainVersionId: mainVersionId2,
			Title:         title2,
			Owners:        owners2,
		},
	}

	// Act
	actual, err := serv.GetArticleList()

	// Assert
	assert.Equal(t, actual, expected)
	assert.Equal(t, err, nil)

	articleRepoMock.Mock.AssertCalled(t, "GetAllArticles", 1)

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": mainVersionId2,
	})
}

func TestGetArticleListDbFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId1 := int64(1)
	articleId2 := int64(2)
	mainVersionId1 := int64(3)
	mainVersionId2 := int64(4)

	repositories.GetAllArticlesMock = func() ([]entities.Article, error) {
		return []entities.Article{
			{
				Id:            articleId1,
				MainVersionID: mainVersionId1,
			}, {
				Id:            articleId2,
				MainVersionID: mainVersionId2,
			},
		}, errors.New("hello I am an error")
	}

	var expected []models.ArticleListElement

	// Act
	actual, err := serv.GetArticleList()

	// Assert
	assert.Equal(t, actual, expected)
	assert.NotEqual(t, err, nil)

	articleRepoMock.Mock.AssertCalled(t, "GetAllArticles", 1)

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 0)
}

func TestGetArticleListElementFail(t *testing.T) {
	// Arrange
	localSetup()

	articleId1 := int64(1)
	articleId2 := int64(2)
	mainVersionId1 := int64(3)
	mainVersionId2 := int64(4)
	title1 := "Lorem ipsum"
	title2 := "Dolor sit amet"
	owners1 := []string{"johndoe@gmail.com"}
	owners2 := []string{"janedoe@gmail.com"}

	repositories.GetAllArticlesMock = func() ([]entities.Article, error) {
		return []entities.Article{
			{
				Id:            articleId1,
				MainVersionID: mainVersionId1,
			}, {
				Id:            articleId2,
				MainVersionID: mainVersionId2,
			},
		}, nil
	}

	repositories.GetVersionMock = func(version int64) (entities.Version, error) {
		if version == mainVersionId1 {
			return entities.Version{
				ArticleID:      articleId1,
				Id:             mainVersionId1,
				Title:          title1,
				Owners:         owners1,
				Status:         "draft",
				LatestCommitID: "commitId",
			}, nil
		} else {
			return entities.Version{
				ArticleID:      articleId2,
				Id:             mainVersionId2,
				Title:          title2,
				Owners:         owners2,
				Status:         "draft",
				LatestCommitID: "commitId",
			}, errors.New("oh no I am an error")
		}
	}

	// although version2 gave an error, the list should still have version1
	expected := []models.ArticleListElement{{
		Id:            articleId1,
		MainVersionId: mainVersionId1,
		Title:         title1,
		Owners:        owners1},
	}

	// Act
	actual, err := serv.GetArticleList()

	// Assert
	assert.Equal(t, actual, expected)
	assert.Equal(t, err, nil)

	articleRepoMock.Mock.AssertCalled(t, "GetAllArticles", 1)

	versionRepoMock.Mock.AssertCalled(t, "GetVersion", 2)
	versionRepoMock.Mock.AssertCalledWith(t, "GetVersion", &map[string]interface{}{
		"version": mainVersionId2,
	})

}

func TestCreateArticleSuccess(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	duplicateOwners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com", "JohnDoe@gmail.com"}
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}
	loggedInAs := "JaneDoe@gmail.com"

	commitId := "thisIsACommitId"
	articleId := int64(1)
	versionId := int64(1)

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateArticleMock = func() (entities.Article, error) {
		return entities.Article{Id: articleId, MainVersionID: versionId}, nil
	}

	repositories.CreateVersionMock = func(version entities.Version) (entities.Version, error) {
		version.Id = articleId
		return version, nil
	}

	repositories.UpdateMainVersionMock = func(id int64, id2 int64) error {
		return nil
	}

	repositories.InitMainVersionMock = func(article int64, mainVersion int64) (string, error) {
		return commitId, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return nil
	}

	expected := models.Version{
		ArticleID:      articleId,
		Id:             versionId,
		Title:          title,
		Owners:         owners,
		Content:        "",
		Status:         "draft",
		LatestCommitID: commitId,
	}

	// Act
	article, err := serv.CreateArticle(title, duplicateOwners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.Equal(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 1)

	versionRepoMock.Mock.AssertCalledWith(t, "CreateVersion", &map[string]interface{}{
		"version": entities.Version{
			ArticleID:      articleId,
			Id:             0,
			Title:          title,
			Owners:         owners,
			Status:         "",
			LatestCommitID: "",
		},
	})

	articleRepoMock.Mock.AssertCalledWith(t, "UpdateMainVersion", &map[string]interface{}{
		"id":  articleId,
		"id2": versionId,
	})

	storerMock.Mock.AssertCalledWith(t, "InitMainVersion", &map[string]interface{}{
		"article":     articleId,
		"mainVersion": versionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "UpdateVersionLatestCommit", &map[string]interface{}{
		"version": versionId,
		"commit":  commitId,
	})
}

func TestCreateArticleNonExistingUserFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	loggedInAs := "JaneDoe@gmail.com"
	notExist := "OhNo@gmail.com"
	owners := []string{loggedInAs, notExist}

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		if email == notExist {
			return false, nil
		}
		return true, nil
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, owners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 0)

	versionRepoMock.Mock.AssertCalled(t, "CreateVersion", 0)

	articleRepoMock.Mock.AssertCalled(t, "UpdateMainVersion", 0)

	storerMock.Mock.AssertCalled(t, "InitMainVersion", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleCannotCheckUserExistsFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	loggedInAs := "JaneDoe@gmail.com"
	owners := []string{loggedInAs, "JohnDoe@gmail.com"}

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, errors.New("haha this is an error")
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, owners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[0],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 0)

	versionRepoMock.Mock.AssertCalled(t, "CreateVersion", 0)

	articleRepoMock.Mock.AssertCalled(t, "UpdateMainVersion", 0)

	storerMock.Mock.AssertCalled(t, "InitMainVersion", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleOwnerNotAuthorizedFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	loggedInAs := "HackerMan@gmail.com"
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, owners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 0)

	versionRepoMock.Mock.AssertCalled(t, "CreateVersion", 0)

	articleRepoMock.Mock.AssertCalled(t, "UpdateMainVersion", 0)

	storerMock.Mock.AssertCalled(t, "InitMainVersion", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleDbArticleCreationFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	duplicateOwners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com", "JohnDoe@gmail.com"}
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}
	loggedInAs := "JaneDoe@gmail.com"

	articleId := int64(1)
	versionId := int64(1)

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateArticleMock = func() (entities.Article, error) {
		return entities.Article{Id: articleId, MainVersionID: versionId},
			errors.New("haha you got an error")
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, duplicateOwners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 1)

	versionRepoMock.Mock.AssertCalled(t, "CreateVersion", 0)

	articleRepoMock.Mock.AssertCalled(t, "UpdateMainVersion", 0)

	storerMock.Mock.AssertCalled(t, "InitMainVersion", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleDbVersionCreationFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	duplicateOwners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com", "JohnDoe@gmail.com"}
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}
	loggedInAs := "JaneDoe@gmail.com"

	articleId := int64(1)
	versionId := int64(1)

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateArticleMock = func() (entities.Article, error) {
		return entities.Article{Id: articleId, MainVersionID: versionId}, nil
	}

	repositories.CreateVersionMock = func(version entities.Version) (entities.Version, error) {
		version.Id = articleId
		return version, errors.New("haha you got an error")
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, duplicateOwners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 1)

	versionRepoMock.Mock.AssertCalledWith(t, "CreateVersion", &map[string]interface{}{
		"version": entities.Version{
			ArticleID:      articleId,
			Id:             0,
			Title:          title,
			Owners:         owners,
			Status:         "",
			LatestCommitID: "",
		},
	})

	articleRepoMock.Mock.AssertCalled(t, "UpdateMainVersion", 0)

	storerMock.Mock.AssertCalled(t, "InitMainVersion", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleUpdateMainVersionFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	duplicateOwners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com", "JohnDoe@gmail.com"}
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}
	loggedInAs := "JaneDoe@gmail.com"

	articleId := int64(1)
	versionId := int64(1)

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateArticleMock = func() (entities.Article, error) {
		return entities.Article{Id: articleId, MainVersionID: versionId}, nil
	}

	repositories.CreateVersionMock = func(version entities.Version) (entities.Version, error) {
		version.Id = articleId
		return version, nil
	}

	repositories.UpdateMainVersionMock = func(id int64, id2 int64) error {
		return errors.New("haha you got an error")
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, duplicateOwners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 1)

	versionRepoMock.Mock.AssertCalledWith(t, "CreateVersion", &map[string]interface{}{
		"version": entities.Version{
			ArticleID:      articleId,
			Id:             0,
			Title:          title,
			Owners:         owners,
			Status:         "",
			LatestCommitID: "",
		},
	})

	articleRepoMock.Mock.AssertCalledWith(t, "UpdateMainVersion", &map[string]interface{}{
		"id":  articleId,
		"id2": versionId,
	})

	storerMock.Mock.AssertCalled(t, "InitMainVersion", 0)

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleInitMainVersionFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	duplicateOwners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com", "JohnDoe@gmail.com"}
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}
	loggedInAs := "JaneDoe@gmail.com"

	commitId := "thisIsACommitID"
	articleId := int64(1)
	versionId := int64(1)

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateArticleMock = func() (entities.Article, error) {
		return entities.Article{Id: articleId, MainVersionID: versionId}, nil
	}

	repositories.CreateVersionMock = func(version entities.Version) (entities.Version, error) {
		version.Id = articleId
		return version, nil
	}

	repositories.UpdateMainVersionMock = func(id int64, id2 int64) error {
		return nil
	}

	repositories.InitMainVersionMock = func(article int64, mainVersion int64) (string, error) {
		return commitId, errors.New("haha you got an error")
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, duplicateOwners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 1)

	versionRepoMock.Mock.AssertCalledWith(t, "CreateVersion", &map[string]interface{}{
		"version": entities.Version{
			ArticleID:      articleId,
			Id:             0,
			Title:          title,
			Owners:         owners,
			Status:         "",
			LatestCommitID: "",
		},
	})

	articleRepoMock.Mock.AssertCalledWith(t, "UpdateMainVersion", &map[string]interface{}{
		"id":  articleId,
		"id2": versionId,
	})

	storerMock.Mock.AssertCalledWith(t, "InitMainVersion", &map[string]interface{}{
		"article":     articleId,
		"mainVersion": versionId,
	})

	versionRepoMock.Mock.AssertCalled(t, "UpdateVersionLatestCommit", 0)
}

func TestCreateArticleUpdateLatestCommitFail(t *testing.T) {
	// Arrange
	localSetup()

	title := "Lorem ipsum dolor sit amet"
	duplicateOwners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com", "JohnDoe@gmail.com"}
	owners := []string{"JaneDoe@gmail.com", "JohnDoe@gmail.com"}
	loggedInAs := "JaneDoe@gmail.com"

	commitId := "thisIsACommitId"
	articleId := int64(1)
	versionId := int64(1)

	// define mock behaviour
	repositories.CheckIfExistsMock = func(email string) (bool, error) {
		return true, nil
	}

	repositories.CreateArticleMock = func() (entities.Article, error) {
		return entities.Article{Id: articleId, MainVersionID: versionId}, nil
	}

	repositories.CreateVersionMock = func(version entities.Version) (entities.Version, error) {
		version.Id = articleId
		return version, nil
	}

	repositories.UpdateMainVersionMock = func(id int64, id2 int64) error {
		return nil
	}

	repositories.InitMainVersionMock = func(article int64, mainVersion int64) (string, error) {
		return commitId, nil
	}

	repositories.UpdateVersionLatestCommitMock = func(version int64, commit string) error {
		return errors.New("haha you got an error")
	}

	expected := models.Version{}

	// Act
	article, err := serv.CreateArticle(title, duplicateOwners, loggedInAs)

	// Assert
	assert.Equal(t, article, expected)
	assert.NotEqual(t, err, nil)

	userRepoMock.Mock.AssertCalledWith(t, "CheckIfExists", &map[string]interface{}{
		"email": owners[1],
	})

	articleRepoMock.Mock.AssertCalled(t, "CreateArticle", 1)

	versionRepoMock.Mock.AssertCalledWith(t, "CreateVersion", &map[string]interface{}{
		"version": entities.Version{
			ArticleID:      articleId,
			Id:             0,
			Title:          title,
			Owners:         owners,
			Status:         "",
			LatestCommitID: "",
		},
	})

	articleRepoMock.Mock.AssertCalledWith(t, "UpdateMainVersion", &map[string]interface{}{
		"id":  articleId,
		"id2": versionId,
	})

	storerMock.Mock.AssertCalledWith(t, "InitMainVersion", &map[string]interface{}{
		"article":     articleId,
		"mainVersion": versionId,
	})

	versionRepoMock.Mock.AssertCalledWith(t, "UpdateVersionLatestCommit", &map[string]interface{}{
		"version": versionId,
		"commit":  commitId,
	})
}
