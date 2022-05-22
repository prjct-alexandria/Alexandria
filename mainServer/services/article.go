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
	article := entities.Article{}
	article, err := serv.articlerepo.SaveArticle(article)
	if err != nil {
		return models.Version{}, err
	}

	// Create article git repo
	err = serv.gitrepo.CreateRepo(article.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Create main version info in database
	version := entities.Version{Title: title, Owners: owners}
	version, err = serv.versionrepo.SaveVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	// Link version to article
	err = serv.articlerepo.LinkVersion(article.Id, version.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Return frontend-readable description of created data
	return models.Version{
		ArticleID: article.Id,
		Id:        version.Id,
		Title:     version.Title,
		Owners:    version.Owners,
		Content:   "",
	}, nil
}
