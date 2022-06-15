package controllers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services/interfaces"
	"mainServer/utils/httperror"
	"net/http"
	"os"
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

	// get optional query parameter for specific history/commit ID
	values := c.Request.URL.Query()
	var usingCommit bool
	var commitHashArr *[20]byte
	if commitStr, ok := values["historyID"]; ok {
		// read the string like 'a8fc73280...' into a byte array for the commit hash
		commitHashSlice, err := hex.DecodeString(commitStr[0])
		if err != nil || len(commitHashSlice) != 20 {
			fmt.Println(err)
			httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("invalid commit id=%s, should be a 40-character long hex string", commitStr[0]))
			return
		}

		// cast the Go slice to a fixed length array, after having checked if the slice had the right length
		commitHashArr = (*[20]byte)(commitHashSlice)
		usingCommit = true
	}

	// Get either a specific version or just the latest
	var res models.Version
	if usingCommit {
		res, err = contr.Serv.GetVersionByCommitID(article, version, *commitHashArr)
	} else {
		res, err = contr.Serv.GetVersion(article, version)
	}

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

	path, err := contr.Serv.GetVersionFiles(article, version)
	if err != nil {
		//TODO create separate error scenarios (article / version doesn't exist, zip failed)
		httperror.NewError(c, http.StatusBadRequest, errors.New("could not get article files"))
		return
	}

	//GetVersionFiles creates a temporary zip file, which needs to be removed after this method is finished

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}(path)

	//Return files
	c.Header("Access-Control-Expose-Headers", "content-disposition")
	c.FileAttachment(path, filepath.Base(path))
}
