package repositories

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

type ArticleRepositoryMock struct {
	Mock *mocks.Mock
}

func NewArticleRepositoryMock() ArticleRepositoryMock {
	return ArticleRepositoryMock{Mock: mocks.NewMock()}
}

var CreateArticleMock func() (entities.Article, error)

func (m ArticleRepositoryMock) CreateArticle() (entities.Article, error) {
	m.Mock.CallFunc("CreateArticle", &map[string]interface{}{})
	return CreateArticleMock()
}

var UpdateMainVersionMock func(id int64, id2 int64) error

func (m ArticleRepositoryMock) UpdateMainVersion(id int64, id2 int64) error {
	m.Mock.CallFunc("UpdateMainVersion", &map[string]interface{}{
		"id":  id,
		"id2": id2,
	})
	return UpdateMainVersionMock(id, id2)
}

var GetMainVersionMock func(article int64) (int64, error)

func (m ArticleRepositoryMock) GetMainVersion(article int64) (int64, error) {
	m.Mock.CallFunc("GetMainVersion", &map[string]interface{}{
		"article": article,
	})
	return GetMainVersionMock(article)
}

var GetAllArticlesMock func() ([]entities.Article, error)

func (m ArticleRepositoryMock) GetAllArticles() ([]entities.Article, error) {
	m.Mock.CallFunc("GetAllArticles", &map[string]interface{}{})
	return GetAllArticlesMock()
}
