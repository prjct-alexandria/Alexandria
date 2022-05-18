package services

import (
	"mainServer/entities"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
)

type ArticleService struct {
	versionrepo interfaces.VersionRepository
	gitrepo     repositories.GitRepository
}

func NewArticleService(versionrepo interfaces.VersionRepository, gitrepo repositories.GitRepository) ArticleService {
	return ArticleService{versionrepo: versionrepo, gitrepo: gitrepo}
}

// CreateArticle creates a new article repo and main article version, returns main version
func (serv ArticleService) CreateArticle(title string, owners []string) (entities.Version, error) {
	// TODO: ensure authenticated user is among owners
	id, err := gitrepo.CreateRepo()

}
