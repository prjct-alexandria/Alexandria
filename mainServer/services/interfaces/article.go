package interfaces

import "mainServer/models"

type ArticleService interface {
	CreateArticle(title string, owners []string) (models.Version, error)
	GetMainVersion(article int64) (int64, error)
	GetArticleList() ([]models.ArticleListElement, error)
}
