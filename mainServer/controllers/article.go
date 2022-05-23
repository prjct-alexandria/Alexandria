package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"mainServer/models"
	"mainServer/services"
	"net/http"
)

type ArticleController struct {
	serv services.ArticleService
}

func NewArticleController(serv services.ArticleService) ArticleController {
	return ArticleController{serv: serv}
}

// CreateArticle godoc
// @Summary      Create new article
// @Description  Creates new article, including main article version. Returns main version. Owners must be specified as email addresses, not usernames.
// @Accept		 json
// @Param		 article 		body	models.ArticleCreationForm		true 	"Article info"
// @Produce      json
// @Success      200  {object} models.Version
// @Router       /articles [post]
func (contr ArticleController) CreateArticle(c *gin.Context) {

	// Read article creation JSON
	article := models.ArticleCreationForm{}
	err := c.ShouldBindBodyWith(article, binding.JSON)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
	}

	// Create article in service
	version, err := contr.serv.CreateArticle(article.Title, article.Owners)

	// Respond with a frontend-readable description of the created version
	c.JSON(http.StatusOK, version)
}
