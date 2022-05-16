package services

import (
	"github.com/gin-gonic/gin"
	"mainServer/repositories"
	"mime/multipart"
	"path/filepath"
)

type VersionService struct {
	Gitrepo repositories.GitRepository
}

// UpdateVersion overwrites file of specified article version and commits
func (serv VersionService) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
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
