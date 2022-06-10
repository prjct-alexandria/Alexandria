package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"mainServer/server/config"
	"mainServer/utils/clock"
	"time"
)

// AuthMiddleware Main source for creating this handler: https://medium.com/wesionary-team/jwt-authentication-in-golang-with-gin-63dbc0816d55
// Middleware function for handling authentication
// It reads the Authorization cookie to get a JWT token
// If no cookie / token is present the context-email is set to nil
// If the token cannot be validated, the context-email is set to nil
// If the token has expired, the context-email is set to nil
// If there is a valid token, the context-email is set to the email of the logged-in user and the cookie is refreshed
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Request.Cookie("Authorization")
		if err != nil {
			c.Set("Email", nil)
			return
		}

		tokenString := authHeader.Value[len("Bearer."):]

		token, err := validateToken(tokenString, []byte(cfg.Auth.JwtSecret))
		if err != nil {
			c.Set("Email", nil)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			//Check expired
			//Should not be executed as the cookie should have been timed out by now, but extra check for security
			exp := time.Unix(int64(claims["expiresAt"].(float64)), 0)
			if (clock.RealClock{}.Now().After(exp)) {
				c.Set("Email", nil)
				//TODO: Notify user on frontend that sesssion has timed out?
				return
			}

			//Store in this context
			c.Set("Email", claims["email"].(string))

			//Refresh token
			err := UpdateJwtCookie(c, claims["email"].(string), cfg)

			if err != nil {
				return
			}
		}
	}
}

// Function for validating JWT token
// Code from https://medium.com/wesionary-team/jwt-authentication-in-golang-with-gin-63dbc0816d55
func validateToken(encodedToken string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("Invalid token")
		}
		return secret, nil
	})
}

// UpdateJwtCookie Function for creating and updating the JWT token
func UpdateJwtCookie(c *gin.Context, email string, cfg *config.Config) error {
	cl := clock.RealClock{}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"expiresAt": cl.Now().Add(time.Minute * time.Duration(cfg.Auth.TokenExpireMinutes)).Unix(),
		"issuedAt":  cl.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.Auth.JwtSecret))

	if err != nil {
		return err
	}

	//TODO: Make secure once HTTPS connection is established
	c.SetCookie("Authorization", "Bearer."+tokenString, 60*15, "/", cfg.Hosting.Backend.Host, false, true)

	return nil
}
