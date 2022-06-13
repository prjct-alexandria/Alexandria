package server

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "mainServer/docs"
	"mainServer/middlewares"
)

// @title API documentation
// @version 1.0.0
// @description Documentation of Alexandria's API. Endpoints can be tried out directly in this interactive documentation.

// @host      localhost:8080

func SetUpRouter(contrs ControllerEnv) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.AuthMiddleware())
	router.Use(middlewares.CorsHeaders())

	router.POST("/articles", contrs.article.CreateArticle)
	router.POST("/articles/:articleID/versions/:versionID", contrs.version.UpdateVersion)
	router.POST("/articles/:articleID/versions", contrs.version.CreateVersionFrom)
	router.GET("/articles/:articleID/versions/:versionID", contrs.version.GetVersion)
	router.GET("/articles/:articleID/mainVersion", contrs.article.GetMainVersion)
	router.GET("/articles/:articleID/versions", contrs.version.ListVersions)

	router.POST("/articles/:articleID/requests", contrs.req.CreateRequest)
	router.PUT("/articles/:articleID/requests/:requestID/reject", contrs.req.RejectRequest)
	router.PUT("/articles/:articleID/requests/:requestID/accept", contrs.req.AcceptRequest)

	router.POST("/users", contrs.user.Register)
	router.POST("/login", contrs.user.Login)
	router.POST("/createExampleUser", contrs.user.CreateExampleUser)

	//Example of how to make an endpoint use the authentication
	router.GET("/getExampleUser", contrs.user.GetExampleUser)

	router.POST("/articles/:articleID/thread/:threadType/id/:specificID/", contrs.thread.CreateThread)
	router.POST("/comments/thread/:threadID", contrs.thread.SaveComment)

	//Groups can be used for nested paths, maybe add example later
	// Path for accessing the API documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
