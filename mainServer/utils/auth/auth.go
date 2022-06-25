package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetLoggedInEmail(c *gin.Context) string {
	email, _ := c.Get("Email")
	if email != nil {
		return fmt.Sprintf("%s", email)
	} else {
		return ""
	}
}

func IsLoggedIn(c *gin.Context) bool {
	_, loggedIn := c.Get("Email")
	return loggedIn
}
