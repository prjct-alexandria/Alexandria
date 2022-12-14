package server

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "mainServer/docs"
	"mainServer/middlewares"
	"mainServer/server/config"
)

// @title API documentation
// @version 1.0.0
// @description Documentation of Alexandria's API. Endpoints can be tried out directly in this interactive documentation.

// @host      localhost:8080

func SetUpRouter(cfg *config.Config, contrs ControllerEnv) *gin.Engine {
	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}

	router.Use(middlewares.AuthMiddleware(cfg))
	router.Use(middlewares.CorsHeaders(cfg))

	router.GET("/articles", contrs.article.ArticleList)
	router.POST("/articles", contrs.article.CreateArticle)
	router.POST("/articles/:articleID/versions/:versionID", contrs.version.UpdateVersion)
	router.POST("/articles/:articleID/versions", contrs.version.CreateVersionFrom)
	router.GET("/articles/:articleID/versions/:versionID", contrs.version.GetVersion)
	router.GET("/articles/:articleID/versions/:versionID/files", contrs.version.GetVersionFiles)
	router.GET("/articles/:articleID/mainVersion", contrs.article.GetMainVersion)
	router.GET("/articles/:articleID/versions", contrs.version.ListVersions)

	router.POST("/articles/:articleID/requests", contrs.req.CreateRequest)
	router.GET("/articles/:articleID/requests", contrs.req.GetRequestList)
	router.GET("/articles/:articleID/requests/:requestID", contrs.req.GetRequest)
	router.PUT("/articles/:articleID/requests/:requestID/reject", contrs.req.RejectRequest)
	router.PUT("/articles/:articleID/requests/:requestID/accept", contrs.req.AcceptRequest)

	router.POST("/users", contrs.user.Register)
	router.POST("/login", contrs.user.Login)
	router.POST("/logout", contrs.user.Logout)

	router.GET("/articles/:articleID/versions/:versionID/history/:commitID/threads",
		contrs.thread.GetCommitThreads)
	router.GET("/articles/:articleID/versions/:versionID/history/:commitID/selectionThreads",
		contrs.thread.GetCommitSelectionThreads)
	router.GET("/articles/:articleID/requests/:requestID/threads", contrs.thread.GetRequestThreads)
	router.POST("/articles/:articleID/thread/:threadType/id/:specificID", contrs.thread.CreateThread)
	router.POST("/comments/thread/:threadID", contrs.thread.SaveComment)

	//Groups can be used for nested paths, maybe add example later
	// Path for accessing the API documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
