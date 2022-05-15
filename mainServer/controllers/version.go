package controllers

import (
	"github.com/gin-gonic/gin"
	"mainServer/repositories"
	"net/http"
	"path/filepath"
)

type VersionController struct {
	gitrepo repositories.GitRepository
}

// UploadFiles godoc
// @Summary     Update article version
// @Description Upload files to update an article version, can only be done by an owner.
// @Accept      mpfd
// @Param		articleID	path	string	true "Article ID"
// @Param		versionID	path	string	true "Version ID"
// @Success     200
// @Failure     400
// @Failure     404
// @Router      /articles/{articleID}/versions/{versionID} [put]
func (contr VersionController) UploadFiles(c *gin.Context) {
	aid := c.Param("articleID")
	vid := c.Param("versionID")

	// Hard coded to only accept /articles/1/versions/1
	// TODO: replace with error if not found in actual database
	if !(aid == "1" && vid == "1") {
		c.Status(http.StatusNotFound)
		return
	}

	// get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, jokes)
}
