package services

import (
	"errors"
	"fmt"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"mainServer/utils/arrays"
)

type ArticleService struct {
	articlerepo interfaces.ArticleRepository
	versionrepo interfaces.VersionRepository
	userrepo    interfaces.UserRepository
	storer      interfaces.Storer
}

func NewArticleService(articlerepo interfaces.ArticleRepository, versionrepo interfaces.VersionRepository, userrepo interfaces.UserRepository, storer interfaces.Storer) ArticleService {
	return ArticleService{
		articlerepo: articlerepo,
		versionrepo: versionrepo,
		userrepo:    userrepo,
		storer:      storer}
}

// CreateArticle creates a new article repo and main article version, returns main version
func (serv ArticleService) CreateArticle(title string, owners []string, loggedInAs string) (models.Version, error) {
	// Remove possible duplicates in owners array
	owners = arrays.RemoveDuplicateStr(owners)

	// Check if owners exist in database
	// Also checks if the authenticated user is in this list
	authenticatedUserPresent := false
	for _, email := range owners {
		exists, err := serv.userrepo.CheckIfExists(email)
		if err != nil {
			return models.Version{}, fmt.Errorf("could not check if %s exists in the database: %s", email, err.Error())
		}
		if !exists {
			return models.Version{}, fmt.Errorf("%s is not a registered email address", email)
		}
		if loggedInAs == email {
			authenticatedUserPresent = true
		}
	}
	// TODO Make this lead to a 403 Forbidden
	if !authenticatedUserPresent {
		return models.Version{}, errors.New("authenticated user is not present in list of owners")
	}

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
		ArticleID:      version.ArticleID,
		Id:             version.Id,
		Title:          version.Title,
		Owners:         version.Owners,
		Content:        "",
		Status:         entities.VersionDraft,
		LatestCommitID: commit,
	}, nil
}

func (serv ArticleService) GetMainVersion(article int64) (int64, error) {
	mv, err := serv.articlerepo.GetMainVersion(article)
	if err != nil {
		return -1, err
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
