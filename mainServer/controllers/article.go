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
// @Param		 title 		body	string		true 	"Title"			example(Lorem Ipsum)
// @Param		 owners	    body	[]string	true	"Owner emails"  example(janedoe@mail.com,johndoe@mail.com)
// @Produce      json
// @Success      200  {object} entities.Article
// @Router       /articles [post]
func (contr ArticleController) CreateArticle(c *gin.Context) {
	// TODO: implement
	c.Status(http.StatusOK)
}
