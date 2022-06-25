package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/services"
	"mainServer/utils/auth"
	"mainServer/utils/httperror"
	"net/http"
	"strconv"
)

type ArticleController struct {
	serv services.ArticleService
}

func NewArticleController(serv services.ArticleService) ArticleController {
	return ArticleController{serv: serv}
}

// CreateArticle godoc
// @Summary      Create new article
// @Description  Creates new article, including main article version. Returns main version info, excluding contents. Owners must be specified as email addresses, not usernames.
// @Accept       json
// @Param        article      body   models.ArticleCreationForm    true   "Article info"
// @Produce      json
// @Success      200  {object} models.Version
// @Failure      400  {object} httperror.HTTPError
// @Failure      500  {object} httperror.HTTPError
// @Router       /articles [post]
func (contr ArticleController) CreateArticle(c *gin.Context) {
	if !auth.IsLoggedIn(c) {
		httperror.NewError(c, http.StatusForbidden, errors.New("must be logged in to perform this request"))
		return
	}
	loggedInAs := auth.GetLoggedInEmail(c)

	// Read article creation JSON
	article := models.ArticleCreationForm{}
	err := c.BindJSON(&article)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, errors.New("cannot bind article creation form"))
		return
	}

	// Create article in service
	version, err := contr.serv.CreateArticle(article.Title, article.Owners, loggedInAs)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusInternalServerError, errors.New("failed creating article on server"))
		return
	}

	// Respond with a frontend-readable description of the created version
	c.JSON(http.StatusOK, version)
}

// GetMainVersion godoc
// @Summary      Get main version id of article
// @Description  Get main version of an article by specifying the article id. Returns the version id of the main version
// @Produce      json
// @Success		 200 "Success"
// @Failure 	 400  {object} httperror.HTTPError
// @Router       /articles/:articleID/mainVersion [get]
func (contr ArticleController) GetMainVersion(c *gin.Context) {
	// extract article id
	aid := c.Param("articleID")
	article, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("Invalid article ID, cannot interpret as integer, id=%s ", aid))
		return
	}

	mv, err := contr.serv.GetMainVersion(article)
	if err != nil {
		fmt.Println(err)
		httperror.NewError(c, http.StatusBadRequest, fmt.Errorf("cannot get main version ID"))
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, strconv.FormatInt(mv, 10))
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
		return
	}
	c.IndentedJSON(http.StatusOK, list)
}
