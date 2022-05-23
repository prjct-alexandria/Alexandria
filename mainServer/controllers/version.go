package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/services"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type VersionController struct {
	Serv services.VersionService
}

// GetVersion 	godoc
// @Summary		Get version content + metadata
// @Description	Gets the version content + metadata from the database + filesystem. Must be accessible without being authenticated.
// @Param		articleID	path	string	true	"Article ID"
// @Param		versionID	path	string	true	"Version ID"
// @Produce		json
// @Success		200 {object} models.Version
// @Failure		404 "Version not found"
// @Router		/articles/{articleID}/versions/{versionID} [get]
func (contr VersionController) GetVersion(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	// extract article id
	aid := c.Param("articleID")
	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// extract version id
	vid := c.Param("versionID")
	version, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid version ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// get version
	res, err := contr.Serv.GetVersion(article, version)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusNotFound, fmt.Errorf("cannot get version with aid=%d and vid=%d", article, version))
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

// UpdateVersion godoc
// @Summary     Update article version
// @Description Upload files to update an article version, can only be done by an owner. Requires multipart form data, with a file attached as the field "file"
// @Accept      mpfd
// @Param		articleID	path	string	true "Article ID"
// @Param		versionID	path	string	true "Version ID"
// @Success     200 "Success"
// @Failure     400 "Bad request, possibly bad file data or permissions"
// @Failure     404 "Specified article version not found"
// @Router      /articles/{articleID}/versions/{versionID} [post]
func (contr VersionController) UpdateVersion(c *gin.Context) {

	// get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// extract article id
	aid := c.Param("articleID")
	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// extract version id
	vid := c.Param("versionID")
	version, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid version ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// update version data
	if err := contr.Serv.UpdateVersion(c, file, article, version); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Status(http.StatusOK)
}
