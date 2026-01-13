package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	// DB dihapus karena tidak digunakan (pakai API Backend)
}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

/* ================================
   SHOW DASHBOARD - Get data dari backend API
================================ */

func (dc *DashboardController) ShowDashboard(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.Redirect(302, "/login")
		return
	}

	// Panggil /barang untuk get all devices
	resp, err := MakeBackendRequest("GET", "/barang", token, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching devices from backend: %v", err)
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":               "Dashboard",
			"totalDevices":        0,
			"activeDevices":       0,
			"locationsCount":      0,
			"alertsCount":         0,
			"activePercentage":    0,
			"deviceLocations":     []gin.H{},
			"deviceLocationsJSON": "[]",
		})
		return
	}

	// Parse response menggunakan helper
	devices, _ := ParseDevicesResponse(resp)
	resp.Body.Close()

	// Hitung statistik
	deviceLocations := []gin.H{}
	totalDevices := len(devices)
	activeDevices := 0
	alertsCount := 0

	for _, device := range devices {
		if device.Status == "Terdeteksi" {
			activeDevices++
		} else {
			alertsCount++
		}

		deviceLocations = append(deviceLocations, gin.H{
			"deviceID": device.DeviceID,
			"name":     device.Name,
			"lat":      device.Latitude,
			"lng":      device.Longitude,
			"status":   device.Status,
		})
	}

	activePercentage := 0
	if totalDevices > 0 {
		activePercentage = (activeDevices * 100) / totalDevices
	}

	// Convert to JSON string untuk digunakan di template
	jsonBytes, _ := json.Marshal(deviceLocations)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":               "Dashboard",
		"totalDevices":        totalDevices,
		"activeDevices":       activeDevices,
		"locationsCount":      activeDevices,
		"alertsCount":         alertsCount,
		"activePercentage":    activePercentage,
		"deviceLocations":     deviceLocations,
		"deviceLocationsJSON": string(jsonBytes),
	})
}

/* ================================
   STREAM DEVICES - SSE endpoint
================================ */

func (dc *DashboardController) StreamDevices(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Set headers untuk SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		log.Println("Streaming not supported")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	ctx := c.Request.Context()

	// Ambil data devices dari backend
	resp, err := MakeBackendRequest("GET", "/barang", token, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching devices: %v", err)
		fmt.Fprintf(c.Writer, "data: []\n\n")
		flusher.Flush()
		return
	}

	// Parse response
	devices, _ := ParseDevicesResponse(resp)
	resp.Body.Close()

	// Transform ke format yang diharapkan frontend
	deviceLocations := []gin.H{}
	for _, device := range devices {
		deviceLocations = append(deviceLocations, gin.H{
			"deviceID": device.DeviceID,
			"name":     device.Name,
			"lat":      device.Latitude,
			"lng":      device.Longitude,
			"status":   device.Status,
		})
	}

	jsonData, _ := json.Marshal(deviceLocations)

	// Kirim ke frontend dengan format SSE
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(jsonData))
	flusher.Flush()

	// Wait untuk client disconnect
	<-ctx.Done()
	log.Println("Client disconnected from SSE")
}
