package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"mainServer/utils/httperror"
	"net/http"
)

//Main source for creating this handler: https://medium.com/wesionary-team/jwt-authentication-in-golang-with-gin-63dbc0816d55

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len("Bearer "):]

		token, err := validateToken(tokenString)
		if err != nil {
			httperror.NewError(c, http.StatusUnauthorized, errors.New("invalid session token"))
			c.Abort()
			return
		}
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			//TODO: add Expire check
			c.Set("Email", claims["email"])
		}

	}
}

func validateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte("temporaryVerySecretThisShouldBeInAConfigFile"), nil
	})
}
