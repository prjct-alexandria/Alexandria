package services

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	Gitrepo     repositories.GitRepository
	Versionrepo interfaces.VersionRepository
}

func (serv VersionService) ListVersions(article int64) ([]models.Version, error) {

	// Get entities from database
	entities, err := serv.Versionrepo.GetVersionsByArticle(article)
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
		Content:   string(fileContent),
		Status:    entity.Status,
	}
	return fullVersion, nil
}

// CreateVersionFrom makes a new version, based of an existing one. Version content is ignored in return value
func (serv VersionService) CreateVersionFrom(article int64, source int64, title string, owners []string) (models.Version, error) {

	// Create entity to store in db
	version := entities.Version{
		ArticleID: article,
		Title:     title,
		Owners:    owners,
	}

	// Store entity in db and receive one with an ID attached
	entity, err := serv.Versionrepo.CreateVersion(version)
	if err != nil {
		return models.Version{}, err
	}

	// Use ID to create new branch in git
	err = serv.Gitrepo.CreateBranch(article, source, entity.Id)
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

func (serv VersionService) GetVersionFiles(article int64, version int64) error {
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

	versionEntity, err := serv.Versionrepo.GetVersion(version)
	versionName := versionEntity.Title

	//TODO: Replace forbidden characters
	versionZip, err := os.Create(filepath.Join(base, "../"+versionName+".zip"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer versionZip.Close()

	zipWriter := zip.NewWriter(versionZip)

	err = addFilesInDirToZip(zipWriter, base, "")
	if err != nil {
		return err
	}
	defer zipWriter.Close()

	//TODO: check if this needs to be moved to the repository
	return nil
}

func addFilesInDirToZip(zipWriter *zip.Writer, dirPath string, dirInZip string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Wrapped in function to allow for "defer file.close()"
		err := func() error {
			if file.IsDir() {
				//Check if it's not the git folder, that's meant for internal server use
				if file.Name() != ".git" {
					err := addFilesInDirToZip(zipWriter, filepath.Join(dirPath, file.Name()), file.Name())
					if err != nil {
						return err
					}
				}
			} else {
				zipFile, err := zipWriter.Create(filepath.Join(dirInZip, file.Name()))
				if err != nil {
					return err
				}

				fileReader, err := os.Open(filepath.Join(dirPath, file.Name()))
				if err != nil {
					return err
				}
				defer fileReader.Close()

				_, err = io.Copy(zipFile, fileReader)
				if err != nil {
					return err
				}
			}
			return nil
		}()

		if err != nil {
			return err
		}
	}
	return nil
}
