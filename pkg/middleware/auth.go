package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/myErrors"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
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
			log.Printf("Middleware: Auth: parse token: %s\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": myErrors.ErrorsString["Unauthorized"]})
		}

		res, ok := token.Claims.(jwt.MapClaims)
		if !ok{
			log.Printf("Middleware: Auth: typecast to jwt.MapClaims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": myErrors.ErrorsString["Unauthorized"]})
		}

		pasCookieRaw := res["password"]
		pasCookie, ok := pasCookieRaw.(string)
		if !ok{
			log.Printf("Middleware: Auth: typecast to string")
			c.JSON(http.StatusUnauthorized, gin.H{"error": myErrors.ErrorsString["Unauthorized"]})
		}


		if pasCookie != pass {
			log.Printf("Middleware: Auth: password is changed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": myErrors.ErrorsString["Unauthorized"]})
		}
	}
	c.Next()
}
