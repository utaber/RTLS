package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// SetupTemplates mengatur template HTML
func SetupTemplates(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
}

// SetupStaticFiles mengatur static files
func SetupStaticFiles(r *gin.Engine) {
	r.Static("/static", "./static")
}

// InitFirebase inisialisasi Firebase Realtime Database
func InitFirebase() *db.Client {
	ctx := context.Background()

	// Path ke service account key
	serviceAccountKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if serviceAccountKey == "" {
		// Coba beberapa lokasi umum
		possiblePaths := []string{
			"service-account-key.json",
			"./service-account-key.json",
			"../service-account-key.json",
			"config/service-account-key.json",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				serviceAccountKey = path
				break
			}
		}

		if serviceAccountKey == "" {
			log.Fatal("Firebase service account key not found. Please check file location.")
		}
	}

	// Database URL
	databaseURL := os.Getenv("FIREBASE_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "https://rtlsrks513-tes-default-rtdb.asia-southeast1.firebasedatabase.app"
	}

	// Initialize Firebase App
	conf := &firebase.Config{
		DatabaseURL: databaseURL,
	}

	opt := option.WithCredentialsFile(serviceAccountKey)
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	// Get Database client
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firebase database: %v\n", err)
	}

	log.Println("Firebase initialized successfully")
	return client
}
