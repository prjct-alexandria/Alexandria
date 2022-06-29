package middlewares

import (
	"github.com/gin-gonic/gin"
	"mainServer/server/config"
	"mainServer/utils/arrays"
	"net/http"
)

//Partially found from: https://github.com/gin-gonic/gin/issues/559#issuecomment-350911039
func CorsHeaders(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := []string{
		cfg.Hosting.Frontend.Url(),
		cfg.Hosting.Frontend.Hostname(),
		cfg.Hosting.Frontend.LocalUrl(),
	}
	return func(c *gin.Context) {
		if arrays.Contains(allowedOrigins, c.GetHeader("origin")) {
			c.Header("Access-Control-Allow-Origin", c.GetHeader("origin"))
		} else {
			// Probably redundant
			c.Header("Access-Control-Allow-Origin", cfg.Hosting.Frontend.LocalUrl())
		}
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}
