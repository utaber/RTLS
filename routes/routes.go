package routes

import (
	"rtls_rks513/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua routes aplikasi
func RegisterRoutes(r *gin.Engine) {
	// Public routes (tidak perlu auth)
	public := r.Group("/")
	{
		public.GET("/", controllers.LoginIndex)
		public.GET("/login", controllers.LoginIndex)
		public.POST("/login", controllers.LoginSubmit)
	}

	// Protected routes (perlu auth)
	protected := r.Group("/")
	// TODO: Uncomment setelah middleware auth dibuat
	// protected.Use(middleware.AuthRequired())
	{
		protected.GET("/dashboard", controllers.DashboardIndex)
		protected.GET("/logout", controllers.Logout)

		// Device UI routes
		protected.GET("/devices", controllers.DevicesIndex)

		// Device REST API routes (untuk modals)
		api := protected.Group("/api/devices")
		{
			api.GET("/:id", controllers.DeviceGet)       // Get single device
			api.POST("", controllers.DeviceCreate)       // Create device
			api.PUT("/:id", controllers.DeviceUpdate)    // Update device
			api.DELETE("/:id", controllers.DeviceDelete) // Delete device
		}

		// TODO: Tambahkan routes lain
		// protected.GET("/history", controllers.HistoryIndex)
	}
}
