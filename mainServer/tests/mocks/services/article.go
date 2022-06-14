package services

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

// ArticleServiceMock mocks the ArticleService with function responses,
// Storing function calls in the Mock field
type ArticleServiceMock struct {
	Mock *mocks.Mock
}

// NewArticleServiceMock initializes a mock that records function calls
func NewArticleServiceMock() ArticleServiceMock {
	return ArticleServiceMock{Mock: mocks.NewMock()}
}

var CreateArticleMock func() (entities.Article, error)

func (m ArticleServiceMock) CreateArticle() (entities.Article, error) {
	m.Mock.CallFunc("CreateArticle", &map[string]interface{}{})
	return CreateArticleMock()
}

var UpdateMainVersionMock func(id int64, id2 int64) error

func (m ArticleServiceMock) UpdateMainVersion(id int64, id2 int64) error {
	m.Mock.CallFunc("UpdateMainVersion", &map[string]interface{}{
		"id":  id,
		"id2": id2,
	})
	return UpdateMainVersionMock(id, id2)
}

var GetMainVersionMock func(article int64) (int64, error)

func (m ArticleServiceMock) GetMainVersion(article int64) (int64, error) {
	m.Mock.CallFunc("GetMainVersion", &map[string]interface{}{
		"article": article,
	})
	return GetMainVersionMock(article)
}

var GetAllArticlesMock func() ([]entities.Article, error)

func (m ArticleServiceMock) GetAllArticles() ([]entities.Article, error) {
	m.Mock.CallFunc("GetAllArticles", &map[string]interface{}{})
	return GetAllArticlesMock()
}
