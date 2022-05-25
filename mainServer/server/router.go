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

	router.POST("/articles/:articleID/versions/:versionID", contrs.version.UpdateVersion)
	router.GET("/articles/:articleID/versions/:versionID", contrs.version.GetVersion)

	router.POST("/users", contrs.user.Register)
	router.POST("/login", contrs.user.Login)

	router.POST("/createExampleUser", contrs.user.CreateExampleUser)

	//Example of how to make an endpoint use the authentication
	router.GET("/getExampleUser", contrs.user.GetExampleUser)

	//Groups can be used for nested paths, maybe add example later
	// Path for accessing the API documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
