package config

import "github.com/gin-gonic/gin"

// SetupTemplates mengatur template HTML
func SetupTemplates(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*.html")
}

// SetupStaticFiles mengatur static files
func SetupStaticFiles(r *gin.Engine) {
	r.Static("/static", "./static")
}
