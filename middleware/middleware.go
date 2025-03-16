package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	token "github.com/mgsquare/go-ecommerce/tokens"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")

		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no authorization header provided"})
			c.Abort()
			return
		}
		token.ValidateToken(ClientToken)
	}
}
