package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services/interfaces"
	"mainServer/utils/auth"
	gitUtils "mainServer/utils/git"
	"mainServer/utils/httperror"
	"net/http"
	"path/filepath"
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
// @Param 		historyID	query	string	false	"History ID"
// @Produce		json
// @Success		200 {object} models.Version
// @Failure 	400 {object} httperror.HTTPError
// @Router		/articles/{articleID}/versions/{versionID} [get]
func (contr VersionController) GetVersion(c *gin.Context) {
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

	// get optional query parameter for specific history/commit ID
	commitID := c.Query("historyID")
	usingCommit := commitID != ""
	if usingCommit && !gitUtils.IsCommitHash(commitID) {
		err := fmt.Errorf("invalid commit id=%s, should be a 40-character long hex string", commitID)
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}

	// Get either a specific version or just the latest
	var res models.Version
	if usingCommit {
		res, err = contr.Serv.GetVersionByCommitID(article, version, commitID)
	} else {
		res, err = contr.Serv.GetVersion(article, version)
	}
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, fmt.Errorf("cannot get version with aid=%d and vid=%d", article, version))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
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

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
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
	// Check if user is logged in
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

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
	if err := contr.Serv.UpdateVersion(c, file, article, version, loggedInAs); err != nil {
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
	// Check if logged in
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

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
	version, err := contr.Serv.CreateVersionFrom(article, form.SourceVersionID, form.Title, form.Owners, loggedInAs)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError,
			errors.New("could not create new version on server, version name might already be in use"))
		return
	}

	// Return version
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, version)
}

// GetVersionFiles 	godoc
// @Summary		Get all the files of a version as a zip
// @Description	Get all the files of an article version as a zip, should be accessible without being authenticated.
// @Param		articleID	path	string	true	"Article ID"
// @Param		versionID	path	string	true	"Version ID"
// @Produce		application/x-zip-compressed
// @Success		200
// @Failure 	400 {object} httperror.HTTPError
// @Router		/articles/{articleID}/versions/{versionID}/files [get]
func (contr VersionController) GetVersionFiles(c *gin.Context) {
	aid := c.Param("articleID")
	vid := c.Param("versionID")

	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("article id must be an integer"))
		return
	}

	version, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("version id must an integer"))
		return
	}

	path, cleanup, err := contr.Serv.GetVersionFiles(article, version)
	defer cleanup() // delete temporary files when done
	if err != nil {
		//TODO create separate error scenarios (article / version doesn't exist, zip failed)
		httperror.NewError(c, http.StatusBadRequest, errors.New("could not get article files"))
		return
	}

	//Return files
	c.Header("Access-Control-Expose-Headers", "content-disposition")
	c.FileAttachment(path, filepath.Base(path))
}
