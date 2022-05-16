package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/services"
	"net/http"
)

type VersionController struct {
	Serv services.VersionService
}

// UpdateVersion godoc
// @Summary     Update article version
// @Description Upload files to update an article version, can only be done by an owner.
// @Accept      mpfd
// @Param		articleID	path	string	true "Article ID"
// @Param		versionID	path	string	true "Version ID"
// @Success     200
// @Failure     400
// @Failure     404
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

	c.Status(http.StatusOK)
}
