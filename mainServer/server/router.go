package server

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "mainServer/docs"
)

// @title API documentation
// @version 1.0.0
// @description Documentation of Alexandria's API. Endpoints can be tried out directly in this interactive documentation.

// @host      localhost:8080

func SetUpRouter(contrs ControllerEnv) *gin.Engine {
	router := gin.Default()

	router.POST("/articles", contrs.article.CreateArticle)
	router.POST("/articles/:articleID/versions/:versionID", contrs.version.UpdateVersion)
	router.GET("/articles/:articleID/versions/:versionID", contrs.version.GetVersion)

	router.POST("/articles/:articleID/requests", contrs.req.CreateRequest)

	router.POST("/createExampleUser", contrs.user.CreateExampleUser)
	router.GET("/getExampleUser", contrs.user.GetExampleUser)

	//Groups can be used for nested paths, maybe add example later
	// Path for accessing the API documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
