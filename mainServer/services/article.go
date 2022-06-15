package services

import (
	"io/ioutil"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"path/filepath"
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

	// Place a default file and create the initial commit
	err = serv.commitDefaultFile(article.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Get the initial commit ID from the git branch
	commit, err := serv.gitrepo.GetLatestCommit(article.Id, version.Id)
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

// commitDefaultFile copies the template file from the resource folder to the given article
// intended for article creation. does not check out any branch.
func (serv ArticleService) commitDefaultFile(article int64) error {

	// Get the path to the repo
	repo, err := serv.gitrepo.GetArticlePath(article)
	if err != nil {
		return err
	}

	// Read the template file
	source, err := filepath.Abs("./resources/defaultArticle.qmd")
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	// Write contents to the main.qmd file in the repo
	target := filepath.Join(repo, "main.qmd")
	err = ioutil.WriteFile(target, input, 0644)
	if err != nil {
		return err
	}

	// Commit the file to the currently checked out branch
	err = serv.gitrepo.Commit(article)
	if err != nil {
		return err
	}
	return nil
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
