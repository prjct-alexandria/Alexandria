package services

import (
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/storer"
)

type ArticleService struct {
	articlerepo interfaces.ArticleRepository
	versionrepo interfaces.VersionRepository
	storer      storer.Storer
}

func NewArticleService(articlerepo interfaces.ArticleRepository, versionrepo interfaces.VersionRepository, storer storer.Storer) ArticleService {
	return ArticleService{
		articlerepo: articlerepo,
		versionrepo: versionrepo,
		storer:      storer}
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

	// Initialize the repository with the specified main version and default article
	commit, err := serv.storer.InitMainVersion(article.Id, version.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Store the commit id in the database
	err = serv.versionrepo.UpdateVersionLatestCommit(version.Id, commit)
	if err != nil {
		return models.Version{}, err
	}

	// Return frontend-readable description of created data, (excluding file contents)
	return models.Version{
		ArticleID: version.ArticleID,
		Id:        version.Id,
		Title:     version.Title,
		Owners:    version.Owners,
		Content:   "",
		Status:    entities.VersionDraft,
	}, nil
}

func (serv ArticleService) GetMainVersion(article int64) (int64, error) {
	mv, err := serv.articlerepo.GetMainVersion(article)
	if err != nil {
		return 0, err
	}
	return mv, nil
}

func (serv ArticleService) GetArticleList() ([]models.ArticleListElement, error) {
	articleList, err := serv.articlerepo.GetAllArticles()

	if err != nil {
		return nil, err
	}

	var res []models.ArticleListElement

	for _, element := range articleList {
		versionData, err := serv.versionrepo.GetVersion(element.MainVersionID)
		if err != nil {
			continue
		}

		listElement := models.ArticleListElement{
			Id:            element.Id,
			MainVersionId: element.MainVersionID,
			Title:         versionData.Title,
			Owners:        versionData.Owners,
			//TODO: Get owners name instead of email?
			//TODO: CreatedAt = ?? (Sort by creation date)
			//TODO: Article Description = ??
		}
		res = append(res, listElement)
	}
	return res, nil
}
