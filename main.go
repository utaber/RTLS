package main

import (
	"html/template"
	"log"
	"rtls_rks513/config"
	"rtls_rks513/routes"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Set mode (bisa pakai environment variable)
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Setup session
	store := cookie.NewStore([]byte("secret-key-change-this-in-production"))
	r.Use(sessions.Sessions("rtls_session", store))

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

	// Initialize Firebase
	dbClient := config.InitFirebase()

	// Register routes dengan Firebase client
	routes.SetupRoutes(r, dbClient)

	// Start server
	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
