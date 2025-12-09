package main

import (
	"html/template"
	"log"
	"rtls_rks513/config"
	"rtls_rks513/routes"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set mode (bisa pakai environment variable)
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Setup templates dengan custom functions
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"year": func() int {
			return time.Now().Year()
		},
	})

	// Setup templates dan static files
	config.SetupTemplates(r)
	config.SetupStaticFiles(r)

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
