package interfaces

import "mainServer/entities"

type ArticleRepository interface {

	// CreateArticle creates a new article and returns the created entity
	// The id in the parameter will be ignored, the generated one is included in the return value
	CreateArticle() (entities.Article, error)

	// LinkVersion links one version to one article
	// An article can link to many versions, but not vice-versa
	LinkVersion(articleID int64, versionID int64) error
}
