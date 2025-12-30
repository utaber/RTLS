package routes

import (
	"rtls_rks513/controllers"
	"rtls_rks513/middleware"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, dbClient *db.Client) {
	// Initialize controllers dengan Firebase client
	dashboardController := controllers.NewDashboardController(dbClient)
	deviceController := controllers.NewDeviceController(dbClient)

	// Public routes (tidak perlu login)
	public := r.Group("/")
	{
		public.GET("/", controllers.LoginIndex)
		public.GET("/login", controllers.LoginIndex)
		public.POST("/login", controllers.LoginSubmit)
		public.GET("/logout", controllers.Logout)
	}

	// Protected routes (HARUS LOGIN)
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Dashboard
		protected.GET("/dashboard", dashboardController.ShowDashboard)

		// SSE endpoint untuk real-time updates
		protected.GET("/stream/devices", dashboardController.StreamDevices)

		// Device pages
		protected.GET("/devices", deviceController.ShowDevices)

		// Device REST API
		api := protected.Group("/api/devices")
		{
			api.GET("/:id", deviceController.DeviceGet)
			api.POST("", deviceController.DeviceCreate)
			api.PUT("/:id", deviceController.DeviceUpdate)
			api.DELETE("/:id", deviceController.DeviceDelete)
		}
	}
}
