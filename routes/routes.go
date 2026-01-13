package routes

import (
	"rtls_rks513/controllers"
	"rtls_rks513/middleware"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, dbClient *db.Client) {
	authCtrl := controllers.NewAuthController()
	dashboardCtrl := controllers.NewDashboardController(dbClient)
	deviceCtrl := controllers.NewDeviceController(dbClient)

	// Public routes
	router.GET("/login", authCtrl.LoginIndex)
	router.POST("/login", authCtrl.LoginSubmit)

	// Protected routes (memerlukan auth)
	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Dashboard
		protected.GET("/dashboard", dashboardCtrl.ShowDashboard)
		protected.GET("/stream/devices", dashboardCtrl.StreamDevices)

		// Devices page
		protected.GET("/devices", deviceCtrl.ShowDevices)

		// Device API endpoints (yang dipanggil oleh JavaScript)
		protected.GET("/api/devices", deviceCtrl.DeviceGet)           //  Ambil device
		protected.POST("/api/devices", deviceCtrl.DeviceCreate)       //  Create device
		protected.PUT("/api/devices/:id", deviceCtrl.DeviceUpdate)    //  Update device
		protected.DELETE("/api/devices/:id", deviceCtrl.DeviceDelete) //  Delete device

		// Logout
		protected.GET("/logout", authCtrl.Logout)
	}
}
