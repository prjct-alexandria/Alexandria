package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services"
	"mainServer/utils/httperror"
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
// @Failure 	 400  {object} httperror.HTTPError
// @Failure 	 500  {object} httperror.HTTPError
// @Router       /articles [post]
func (contr ArticleController) CreateArticle(c *gin.Context) {

	// Read article creation JSON
	article := models.ArticleCreationForm{}
	err := c.BindJSON(&article)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, errors.New("cannot bind article creation form"))
		return
	}

	// Create article in service
	version, err := contr.serv.CreateArticle(article.Title, article.Owners)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed creating article on server"))
		return
	}

	// Respond with a frontend-readable description of the created version
	c.JSON(http.StatusOK, version)
}

// ArticleList godoc
// @Summary      Get a list of all articles
// @Description  Gets a list of all articles in the database + some metadata about the main version.
// @Produce      json
// @Success      200  {object} []models.ArticleListElement
// @Failure 	 500  "server could not retrieve article list"
// @Router       /articles [get]
func (contr ArticleController) ArticleList(c *gin.Context) {
	list, err := contr.serv.GetArticleList()
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("server could not retrieve article list"))
	}
	c.IndentedJSON(http.StatusOK, list)
}
