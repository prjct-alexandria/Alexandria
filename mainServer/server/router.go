package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"mainServer/controllers"
	_ "mainServer/docs"
	"mainServer/repositories"
	"mainServer/services"
)

// @title API documentation
// @version 1.0.0
// @description Documentation of Alexandria's API. Endpoints can be tried out directly in this interactive documentation.

// @host      localhost:8080

func SetUpRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	//router.Use(gin.Recovery())
	//router.Use(gin.Logger())

	helloWorldController := new(controllers.HelloWorldController)
	router.GET("/helloWorldJson", helloWorldController.GetHelloWorldJson)

	userController := controllers.UserController{UserService: services.UserService{UserRepository: repositories.NewPgUserRepository(db)}}

	router.POST("/createExampleUser", userController.CreateExampleUser)
	router.GET("/getExampleUser", userController.GetExampleUser)

	//Groups can be used for nested paths, maybe add example later
	// Path for accessing the API documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
