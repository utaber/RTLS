package controllers

import (
	"log"
	"net/http"
	"rtls_rks513/models"

	"github.com/gin-gonic/gin"
)

// DashboardIndex menampilkan halaman dashboard utama
func DashboardIndex(c *gin.Context) {
	// Ambil semua devices dari FastAPI
	devices, err := models.GetAllDevices()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		// Tampilkan dashboard dengan data kosong jika error
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":           "Dashboard",
			"activeDevices":   0,
			"locationsCount":  0,
			"alertsCount":     0,
			"totalDevices":    0,
			"deviceLocations": []map[string]interface{}{},
			"error":           "Failed to load devices from server",
		})
		return
	}

	// Hitung active devices (status: Terdeteksi)
	activeDevices := 0
	for _, device := range devices {
		if device.Status == "Terdeteksi" {
			activeDevices++
		}
	}

	// Hitung unique locations berdasarkan koordinat
	locationMap := make(map[string]bool)
	for _, device := range devices {
		coordKey := device.DeviceID // Gunakan DeviceID sebagai unique identifier
		locationMap[coordKey] = true
	}
	locationsCount := len(locationMap)

	// Hitung alerts (devices yang Tidak_Terdeteksi)
	alertsCount := 0
	for _, device := range devices {
		if device.Status == "Tidak_Terdeteksi" {
			alertsCount++
		}
	}

	// Data untuk Maps - langsung ambil dari device
	deviceLocations := make([]map[string]interface{}, 0)
	for _, device := range devices {
		deviceLocations = append(deviceLocations, map[string]interface{}{
			"name":     device.Name,
			"deviceID": device.DeviceID,
			"status":   device.Status,
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
