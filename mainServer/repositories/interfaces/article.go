package interfaces

import "mainServer/entities"

type ArticleRepository interface {

	// CreateArticle creates a new article and returns the created entity
	// The id in the parameter will be ignored, the generated one is included in the return value
	CreateArticle() (entities.Article, error)
}
