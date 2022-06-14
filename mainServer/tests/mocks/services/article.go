package services

import "mainServer/entities"

// ArticleServiceMock mocks class using publicly modifiable mock functions
type ArticleServiceMock struct {
	// mock tracks what functions were called and with what parameters
	Called *map[string]bool
	Params *map[string]map[string]interface{}
}

// NewArticleServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewArticleServiceMock() ArticleServiceMock {
	return ArticleServiceMock{
		Called: &map[string]bool{},
		Params: &map[string]map[string]interface{}{},
	}
}

var CreateArticleMock func() (entities.Article, error)
var UpdateMainVersionMock func(id int64, id2 int64) error
var GetMainVersionMock func(article int64) (int64, error)
var GetAllArticlesMock func() ([]entities.Article, error)

func (m ArticleServiceMock) CreateArticle() (entities.Article, error) {
	(*m.Called)["CreateArticle"] = true
	return CreateArticleMock()
}

func (m ArticleServiceMock) UpdateMainVersion(id int64, id2 int64) error {
	(*m.Called)["UpdateMainVersion"] = true
	(*m.Params)["UpdateMainVersion"] = map[string]interface{}{
		"id":  id,
		"id2": id2,
	}
	return UpdateMainVersionMock(id, id2)
}

func (m ArticleServiceMock) GetMainVersion(article int64) (int64, error) {
	(*m.Called)["GetMainVersion"] = true
	(*m.Params)["GetMainVersion"] = map[string]interface{}{
		"article": article,
	}
	return GetMainVersionMock(article)
}

func (m ArticleServiceMock) GetAllArticles() ([]entities.Article, error) {
	(*m.Called)["GetAllArticles"] = true
	return GetAllArticlesMock()
}
