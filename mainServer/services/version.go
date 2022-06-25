package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/storer"
	"mime/multipart"
)

type VersionService struct {
	VersionRepo interfaces.VersionRepository
	Storer      storer.Storer
	UserRepo    interfaces.UserRepository
}

func (serv VersionService) ListVersions(article int64) ([]models.Version, error) {

	// Get list from database
	list, err := serv.VersionRepo.GetVersionsByArticle(article)
	if err != nil {
		return nil, err
	}

	// Convert list to models to send to frontend
	result := make([]models.Version, len(list))
	for i, e := range list {
		result[i] = models.Version{
			ArticleID:      article,
			Id:             e.Id,
			Title:          e.Title,
			Owners:         e.Owners,
			Content:        "",
			Status:         e.Status,
			LatestCommitID: e.LatestCommitID,
		}
	}
	return result, nil
}

// GetVersion looks for a version in the filesystem and creates a version entity from it,
// with the appropriate metadata and contents.
func (serv VersionService) GetVersion(article int64, version int64) (models.Version, error) {

	// Get file contents from the (git) file system
	content, err := serv.Storer.GetVersion(article, version)
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
		Content:        content,
		Status:         entity.Status,
		LatestCommitID: entity.LatestCommitID,
	}
	return fullVersion, nil
}

func (serv VersionService) GetVersionByCommitID(article int64, version int64, commit string) (models.Version, error) {

	// TODO: check if the commit is actually part of the specified version,
	// can be done once the version-commit table exists

	// Get file contents from the (git) file system
	content, err := serv.Storer.GetVersionByCommit(article, commit)
	if err != nil {
		return models.Version{}, err
	}

	// Get other version info from database, this should be the same for every commit
	entity, err := serv.VersionRepo.GetVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	fullVersion := models.Version{
		ArticleID:      entity.ArticleID,
		Id:             entity.Id,
		Title:          entity.Title,
		Owners:         entity.Owners,
		Content:        content,
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
	created, err := serv.VersionRepo.CreateVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	// Use ID to create new branch in git with that name
	commit, err := serv.Storer.CreateVersionFrom(article, source, created.Id)
	if err != nil {
		return models.Version{}, err
	}

	// Store the commit id in the database
	err = serv.VersionRepo.UpdateVersionLatestCommit(created.Id, commit)
	if err != nil {
		return models.Version{}, err
	}

	// Return model, made from the new created version entity
	return models.Version{
		ArticleID:      created.ArticleID,
		Id:             created.Id,
		Title:          created.Title,
		Owners:         created.Owners,
		Content:        "",
		LatestCommitID: commit,
	}, nil
}

// UpdateVersion overwrites file of specified article version and commits
func (serv VersionService) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {

	// Update the version contents in the (git) file system
	commit, err := serv.Storer.UpdateAndCommit(c, file, article, version)
	if err != nil {
		return err
	}

	// Store the commit id in the database
	return serv.VersionRepo.UpdateVersionLatestCommit(version, commit)
}

func (serv VersionService) GetVersionFiles(article int64, version int64) (string, func(), error) {

	// Get information about the version from the database
	versionEntity, err := serv.VersionRepo.GetVersion(version)
	if err != nil {
		return "", nil, nil
	}

	// Get a path to the version file contents, zipped
	versionName := versionEntity.Title
	return serv.Storer.GetVersionZipped(article, version, versionName)
}
