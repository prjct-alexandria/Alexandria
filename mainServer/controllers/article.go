package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleController struct{}

// CreateArticle godoc
// @Summary      Create new article
// @Description  Creates new article, including main article version
// @Accept		 json
// @Param		 article 		body	models.ArticleCreationForm		true 	"Article"			example(Lorem Ipsum)
// @Produce      json
// @Success      200  {object} models.Article
// @Router       /articles [post]
func (contr ArticleController) CreateArticle(c *gin.Context) {
	// TODO: implement
	c.Status(http.StatusOK)
}
