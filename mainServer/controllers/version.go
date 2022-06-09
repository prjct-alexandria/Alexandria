package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services/interfaces"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type VersionController struct {
	Serv interfaces.VersionService
}

// GetVersion 	godoc
// @Summary		Get version content + metadata
// @Description	Gets the version content + metadata from the database + filesystem. Must be accessible without being authenticated.
// @Param		articleID	path	string	true	"Article ID"
// @Param		versionID	path	string	true	"Version ID"
// @Produce		json
// @Success		200 {object} models.Version
// @Failure 	400 {object} httperror.HTTPError
// @Router		/articles/{articleID}/versions/{versionID} [get]
func (contr VersionController) GetVersion(c *gin.Context) {
	c.Header("Content-Type", "application/json")

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

// ListVersions 	godoc
// @Summary		List article versions
// @Description	Gets all versions belonging to a specific article. Does not include version contents.
// @Param		articleID	path	string	true	"Article ID"
// @Produce		json
// @Success		200 {object} []models.Version
// @Failure 	400  {object} httperror.HTTPError
// @Failure 	500  {object} httperror.HTTPError
// @Router		/articles/{articleID}/versions [get]
func (contr VersionController) ListVersions(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// extract article id
	aid := c.Param("articleID")
	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	// get versions
	res, err := contr.Serv.ListVersions(article)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, fmt.Errorf("cannot get versions of aid=%d", article))
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
// @Failure 	400  {object} httperror.HTTPError
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
	c.Status(http.StatusOK)
}

// CreateVersionFrom godoc
// @Summary      Create new version
// @Description  Creates new version from an existing one of the same article
// @Accept		 json
// @Param		 articleID		path	string							true 	"Article ID"
// @Param		 version 		body	models.VersionCreationForm		true 	"Version info"
// @Produce      json
// @Success      200  {object} models.Version
// @Failure      400  {object} httperror.HTTPError
// @Failure      500  {object} httperror.HTTPError
// @Router       /articles/{articleID}/versions [post]
func (contr VersionController) CreateVersionFrom(c *gin.Context) {

	// Extract article id
	aid := c.Param("articleID")
	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("invalid article ID, cannot interpret as integer, id=%s", aid))
		return
	}

	// Read version creation JSON
	form := models.VersionCreationForm{}
	err = c.BindJSON(&form)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, errors.New("cannot bind version creation form"))
		return
	}

	// Create version
	version, err := contr.Serv.CreateVersionFrom(article, form.SourceVersionID, form.Title, form.Owners)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("could not create new version on server"))
		return
	}

	// Return version
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, version)
}

func (contr VersionController) GetVersionFiles(c *gin.Context) {
	aid := c.Param("articleID")
	vid := c.Param("versionID")

	article, err := strconv.ParseInt(aid, 10, 64)
	version, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("article and version id must be integers"))
		return
	}

	_, err = contr.Serv.GetVersionFiles(article, version)
	if err != nil {
		//TODO create separate error scenarios (article / version doesn't exist, zip failed)
		httperror.NewError(c, http.StatusBadRequest, errors.New("could not get article files"))
		return
	}

	//GetVersionFiles creates a temporary zip file, which needs to be removed after this method is finished
	//TODO clear cache/downloads
	//Return files
}
