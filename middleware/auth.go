package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired mengecek apakah user sudah login
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari cookie
		token, err := c.Cookie("auth_token")

		// Jika tidak ada token di cookie, coba dari header Authorization
		if err != nil || token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Jika masih tidak ada token
		if token == "" {
			// Untuk API request, return JSON error
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized - No token provided",
				})
				c.Abort()
				return
			}

			// Untuk web request, redirect ke login
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Simpan token ke context untuk digunakan di handler
		c.Set("auth_token", token)

		// Lanjut ke handler berikutnya
		c.Next()
	}
}
