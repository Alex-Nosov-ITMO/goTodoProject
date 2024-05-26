package middleware

import (
	"net/http"
	"fmt"
	//	"strings"
	"os"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *gin.Context) {
	pass := os.Getenv("TODO_PASSWORD")

	if len(pass) > 0 {

		var tokenStr string

		cookie, err := c.Cookie("token")
		if err == nil {
			tokenStr = cookie
		}

	

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return structures.Secret, nil
		})

		if err != nil {
			errStr := fmt.Sprintf("failed to parse token: %s\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errStr})
		}

		res, ok := token.Claims.(jwt.MapClaims)
		if !ok{
			errStr := fmt.Sprintf("failed to parse token: %s\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errStr})
		}

		pasCookieRaw := res["password"]
		pasCookie, ok := pasCookieRaw.(string)
		if !ok{
			errStr := fmt.Sprintf("failed to typecast to string: %s\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errStr})
		}


		if pasCookie != pass {
			errStr := fmt.Sprintf("password is changed: %s\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errStr})
		}

	}
	c.Next()
}
