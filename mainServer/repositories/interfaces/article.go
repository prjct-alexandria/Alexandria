package interfaces

import "mainServer/entities"

type ArticleRepository interface {

	// CreateArticle creates a new article and returns the created entity
	// The id in the parameter will be ignored, the generated one is included in the return value
	CreateArticle() (entities.Article, error)

	// UpdateMainVersion specifies the main version that belongs to an article,
	// this should only be set during the initial article creation step
	UpdateMainVersion(id int64, id2 int64) error

	GetMainVersion(article int64) (int64, error)

	GetAllArticles() ([]entities.Article, error)
}
