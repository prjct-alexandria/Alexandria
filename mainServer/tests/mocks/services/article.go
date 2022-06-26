package services

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type ArticleServiceMock struct {
	Mock *mocks.Mock
}

func NewArticleServiceMock() ArticleServiceMock {
	return ArticleServiceMock{Mock: mocks.NewMock()}
}

var CreateArticleMock func(title string, owners []string, loggedInAs string) (models.Version, error)

func (m ArticleServiceMock) CreateArticle(title string, owners []string, loggedInAs string) (models.Version, error) {
	m.Mock.CallFunc("CreateArticle", &map[string]interface{}{
		"title":      title,
		"owners":     owners,
		"loggedInAs": loggedInAs,
	})
	return CreateArticleMock(title, owners, loggedInAs)
}

var GetMainVersionMock func(article int64) (int64, error)

func (m ArticleServiceMock) GetMainVersion(article int64) (int64, error) {
	m.Mock.CallFunc("GetMainVersion", &map[string]interface{}{
		"article": article,
	})
	return GetMainVersionMock(article)
}

var GetArticleListMock func() ([]models.ArticleListElement, error)

func (m ArticleServiceMock) GetArticleList() ([]models.ArticleListElement, error) {
	m.Mock.CallFunc("GetArticleList", &map[string]interface{}{})
	return GetArticleListMock()
}
