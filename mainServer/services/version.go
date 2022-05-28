package services

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mainServer/models"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"mime/multipart"
	"path/filepath"
)

type VersionService struct {
	Gitrepo     repositories.GitRepository
	Versionrepo interfaces.VersionRepository
}

// GetVersion looks for a version in the filesystem and creates a version entity from it with the appropriate metadata.
func (serv VersionService) GetVersion(article int64, version int64) (models.Version, error) {

	// Get file contents from Git
	err := serv.Gitrepo.CheckoutBranch(article, version)
	if err != nil {
		return models.Version{}, err
	}

	path, err := serv.Gitrepo.GetArticlePath(article)
	if err != nil {
		return models.Version{}, err
	}

	fileContent, err := ioutil.ReadFile(filepath.Join(path, "main.qmd"))
	if err != nil {
		return models.Version{}, err
	}

	// Get other version info from database
	entity, err := serv.Versionrepo.GetVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	fullVersion := models.Version{
		ArticleID: entity.ArticleID,
		Id:        entity.Id,
		Title:     entity.Title,
		Owners:    entity.Owners,
		Content:   string(fileContent)}
	return fullVersion, err
}

// UpdateVersion overwrites file of specified article version and commits
func (serv VersionService) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {
	// TODO: check if user of authenticated session is version owner

	// Checkout
	err := serv.Gitrepo.CheckoutBranch(article, version)
	if err != nil {
		return err
	}

	// Get folder to save file to
	base, err := serv.Gitrepo.GetArticlePath(article)
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
	err = serv.Gitrepo.Commit(article)
	if err != nil {
		return err
	}
	return nil
}
