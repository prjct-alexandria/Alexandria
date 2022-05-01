package server

import (
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
)

func SetUpRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	helloWorldController := new(controllers.HelloWorldController)
	router.GET("/helloWorld", helloWorldController.GetHelloWorld)

	//Groups can be used for nested paths, maybe add example later

	return router
}
