package controllers

import (
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
// @Router      /articles/{articleID}/versions/{versionID} [put]
func (contr VersionController) UpdateVersion(c *gin.Context) {
	// get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	aid := c.Param("articleID")
	vid := c.Param("versionID")

	if err := contr.Serv.UpdateVersion(c, file, aid, vid); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, jokes)
}
