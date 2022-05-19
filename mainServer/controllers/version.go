package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mainServer/services"
	"net/http"
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
// @Success		200 {object} []entities.Version
// @Failure		404 "Version not found"
// @Router		/articles/{articleID}/versions/{versionID} [get]
func (contr VersionController) GetVersion(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	aid := c.Param("articleID")
	vid := c.Param("versionID")

	res, err := contr.Serv.GetVersion(c, aid, vid)
	if err != nil {
		c.Status(http.StatusNotFound)
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, [1]entities.Version{res})
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

	aid := c.Param("articleID")
	vid := c.Param("versionID")

	if err := contr.Serv.UpdateVersion(c, file, aid, vid); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Status(http.StatusOK)
}
