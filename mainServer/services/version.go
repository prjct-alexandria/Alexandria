package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"mime/multipart"
	"os"
	"path/filepath"
)

type VersionService struct {
	GitRepo        repositories.GitRepository
	VersionRepo    interfaces.VersionRepository
	UserRepo       interfaces.UserRepository
	FilesystemRepo repositories.FilesystemRepository
}

func (serv VersionService) IsVersionOwner(article int64, version int64, email string) bool {
	//return serv.VersionRepo.CheckIfOwner()
	return false
}

func (serv VersionService) GetVersionByCommitID(article int64, version int64, commit [20]byte) (models.Version, error) {

	// Get file contents from Git
	err := serv.GitRepo.CheckoutCommit(article, commit)
	if err != nil {
		return models.Version{}, err
	}

	path, err := serv.GitRepo.GetArticlePath(article)
	if err != nil {
		return models.Version{}, err
	}

	fileContent, err := ioutil.ReadFile(filepath.Join(path, "main.qmd"))
	if err != nil {
		return models.Version{}, err
	}

	// Get other version info from database
	entity, err := serv.VersionRepo.GetVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	fullVersion := models.Version{
		ArticleID:      entity.ArticleID,
		Id:             entity.Id,
		Title:          entity.Title,
		Owners:         entity.Owners,
		Content:        string(fileContent),
		Status:         entity.Status,
		LatestCommitID: entity.LatestCommitID,
	}
	return fullVersion, nil
}

func (serv VersionService) ListVersions(article int64) ([]models.Version, error) {

	// Get entities from database
	entities, err := serv.VersionRepo.GetVersionsByArticle(article)
	if err != nil {
		return nil, err
	}

	// Convert entities to models to send to frontend
	result := make([]models.Version, len(entities))
	for i, e := range entities {
		result[i] = models.Version{
			ArticleID: article,
			Id:        e.Id,
			Title:     e.Title,
			Owners:    e.Owners,
			Content:   "",
			Status:    e.Status,
		}
	}
	return result, nil
}

// GetVersion looks for a version in the filesystem and creates a version entity from it with the appropriate metadata.
func (serv VersionService) GetVersion(article int64, version int64) (models.Version, error) {

	// Get file contents from Git
	err := serv.GitRepo.CheckoutBranch(article, version)
	if err != nil {
		return models.Version{}, err
	}

	path, err := serv.GitRepo.GetArticlePath(article)
	if err != nil {
		return models.Version{}, err
	}

	fileContent, err := ioutil.ReadFile(filepath.Join(path, "main.qmd"))
	if err != nil {
		return models.Version{}, err
	}

	// Get other version info from database
	entity, err := serv.VersionRepo.GetVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	fullVersion := models.Version{
		ArticleID:      entity.ArticleID,
		Id:             entity.Id,
		Title:          entity.Title,
		Owners:         entity.Owners,
		Content:        string(fileContent),
		Status:         entity.Status,
		LatestCommitID: entity.LatestCommitID,
	}
	return fullVersion, nil
}

// CreateVersionFrom makes a new version, based of an existing one. Version content is ignored in return value
func (serv VersionService) CreateVersionFrom(article int64, source int64, title string, owners []string) (models.Version, error) {
	// Check if owners exist in database
	for _, email := range owners {
		exists, err := serv.UserRepo.CheckIfExists(email)
		if err != nil {
			return models.Version{}, errors.New(fmt.Sprintf("could not check if %s exists in the database: %s", email, err.Error()))
		}
		if !exists {
			return models.Version{}, errors.New(fmt.Sprintf("%s is not a registered email address", email))
		}
	}

	// Create entity to store in db
	version := entities.Version{
		ArticleID: article,
		Title:     title,
		Owners:    owners,
	}

	// Store entity in db and receive one with an ID attached
	entity, err := serv.VersionRepo.CreateVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	// Use ID to create new branch in git
	err = serv.GitRepo.CreateBranch(article, source, entity.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Store the latest git commit ID of the version in the database entity
	err = serv.UpdateLatestCommit(article, entity.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Return model, made from entity
	return models.Version{
		ArticleID: entity.ArticleID,
		Id:        entity.Id,
		Title:     entity.Title,
		Owners:    entity.Owners,
		Content:   "",
	}, nil
}

// UpdateVersion overwrites file of specified article version and commits
func (serv VersionService) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
	// Check if owner of version

	// Checkout
	err := serv.GitRepo.CheckoutBranch(article, version)
	if err != nil {
		return err
	}

	// Get folder to save file to
	base, err := serv.GitRepo.GetArticlePath(article)
	if err != nil {
		return err
	}

	// Save file
	// TODO: find something more flexible than hard-coding main.qmd
	path := filepath.Join(base, "main.qmd")
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		return err
	}

	// Commit
	err = serv.GitRepo.Commit(article)
	if err != nil {
		return err
	}

	// Store the latest git commit ID of the version in the database entity
	err = serv.UpdateLatestCommit(article, version)
	if err != nil {
		return err
	}

	return nil
}

// UpdateLatestCommit stores the latest git commit ID of the version in the database entity
func (serv VersionService) UpdateLatestCommit(article int64, version int64) error {
	// Get the latest commit ID from the git branch
	commit, err := serv.GitRepo.GetLatestCommit(article, version)
	if err != nil {
		return err
	}
	// Store the commit id in the database
	err = serv.VersionRepo.UpdateVersionLatestCommit(version, commit)
	if err != nil {
		return err
	}
	return nil
}

func (serv VersionService) GetVersionFiles(article int64, version int64) (string, error) {
	// Checkout
	err := serv.GitRepo.CheckoutBranch(article, version)
	if err != nil {
		return "", err
	}

	// Get folder to save file to
	base, err := serv.GitRepo.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	versionEntity, err := serv.VersionRepo.GetVersion(version)
	versionName := versionEntity.Title

	path := filepath.Join(serv.GitRepo.Path, "cache", "downloads", versionName+".zip")
	versionZip, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return path, err
	}
	defer versionZip.Close()

	zipWriter := zip.NewWriter(versionZip)

	err = serv.FilesystemRepo.AddFilesInDirToZip(zipWriter, base, "")
	if err != nil {
		return path, err
	}

	defer zipWriter.Close()

	return path, nil
}
