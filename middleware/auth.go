package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRequired mengecek apakah user sudah login
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ambil cookie token
		token, err := c.Cookie("firebase_token")
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// simpan token ke context (opsional)
		c.Set("firebase_token", token)

		c.Next()
	}
}
