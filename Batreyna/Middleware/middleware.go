package Middleware

import (
	"Batreyna/Token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := Token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized Token Invalid")
			c.Abort()
			return
		}
		c.Next()
	}
}
