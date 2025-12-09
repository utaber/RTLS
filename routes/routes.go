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

		// Device routes (RESTful style)
		protected.GET("/devices", controllers.DevicesIndex)
		protected.GET("/devices/add", controllers.DeviceAddForm)
		protected.POST("/devices/add", controllers.DeviceAdd)
		protected.GET("/devices/:id/edit", controllers.DeviceEditForm)
		protected.POST("/devices/:id/edit", controllers.DeviceEdit)
		protected.GET("/devices/:id/delete", controllers.DeviceDeleteForm)
		protected.POST("/devices/:id/delete", controllers.DeviceDelete)

		// API endpoint untuk modal (get device data)
		protected.GET("/api/devices/:id", controllers.DeviceGet)

	}
}
