package server

import (
	"github.com/gin-gonic/gin"
	"mainServer/controllers"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	//router.Use(gin.Recovery())
	//router.Use(gin.Logger())

	helloWorldController := new(controllers.HelloWorldController)
	router.GET("/helloWorldJson", helloWorldController.GetHelloWorldJson)

	//Groups can be used for nested paths, maybe add example later

	return router
}
