package services

import (
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
)

type ArticleService struct {
	articlerepo interfaces.ArticleRepository
	versionrepo interfaces.VersionRepository
	gitrepo     repositories.GitRepository
}

func NewArticleService(articlerepo interfaces.ArticleRepository, versionrepo interfaces.VersionRepository, gitrepo repositories.GitRepository) ArticleService {
	return ArticleService{
		articlerepo: articlerepo,
		versionrepo: versionrepo,
		gitrepo:     gitrepo}
}

// CreateArticle creates a new article repo and main article version, returns main version
func (serv ArticleService) CreateArticle(title string, owners []string) (models.Version, error) {
	// TODO: ensure authenticated user is among owners

	// Create article in database, this generates article ID
	article, err := serv.articlerepo.CreateArticle()
	if err != nil {
		return models.Version{}, err
	}

	// Create main version in database, this generates version ID
	version := entities.Version{ArticleID: article.Id, Title: title, Owners: owners}
	version, err = serv.versionrepo.CreateVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	// Go back to the article entity and link the main version
	err = serv.articlerepo.UpdateMainVersion(article.Id, version.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Create article git repo
	err = serv.gitrepo.CreateRepo(article.Id, version.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Return frontend-readable description of created data
	return models.Version{
		ArticleID: version.ArticleID,
		Id:        version.Id,
		Title:     version.Title,
		Owners:    version.Owners,
		Content:   "",
	}, nil
}
