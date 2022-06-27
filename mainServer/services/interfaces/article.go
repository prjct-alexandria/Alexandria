package interfaces

import "mainServer/models"

type ArticleService interface {

	// CreateArticle creates an article with the specified title and owners
	// Returns the main version that is automatically created
	CreateArticle(title string, owners []string, loggedInAs string) (models.Version, error)

	// GetMainVersion returns the ID of the main version of the article
	GetMainVersion(article int64) (int64, error)

	// GetArticleList returns a list of all the articles in the database
	GetArticleList() ([]models.ArticleListElement, error)
}
