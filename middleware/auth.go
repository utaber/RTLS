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

		// TODO: Tambahkan routes lain
		protected.GET("/devices", controllers.DevicesIndex)
	}
}
