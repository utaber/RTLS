package controllers

import (
	"net/http"
	"rtls_rks513/models"

	"github.com/gin-gonic/gin"
)

// DashboardIndex menampilkan halaman dashboard utama
func DashboardIndex(c *gin.Context) {
	// Ambil semua devices
	devices := models.GetAllDevices()

	// Hitung active devices
	activeDevices := 0
	for _, device := range devices {
		if device.Status == "active" {
			activeDevices++
		}
	}

	// Hitung unique locations
	locationMap := make(map[string]bool)
	for _, device := range devices {
		locationMap[device.Location] = true
	}
	locationsCount := len(locationMap)

	// Hitung alerts (devices yang inactive)
	alertsCount := 0
	for _, device := range devices {
		if device.Status == "inactive" {
			alertsCount++
		}
	}

	// Data untuk Maps - langsung ambil dari device
	deviceLocations := make([]map[string]interface{}, 0)
	for _, device := range devices {
		deviceLocations = append(deviceLocations, map[string]interface{}{
			"name":     device.Name,
			"deviceID": device.DeviceID,
			"location": device.Location,
			"status":   device.Status,
			"type":     device.Type,
			"lat":      device.Latitude,
			"lng":      device.Longitude,
		})
	}

	data := gin.H{
		"title":           "Dashboard",
		"activeDevices":   activeDevices,
		"locationsCount":  locationsCount,
		"alertsCount":     alertsCount,
		"totalDevices":    len(devices),
		"deviceLocations": deviceLocations,
	}

	c.HTML(http.StatusOK, "dashboard.html", data)
}
