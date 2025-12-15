package controllers

import (
	"log"
	"net/http"
	"rtls_rks513/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DevicesIndex menampilkan daftar semua devices
func DevicesIndex(c *gin.Context) {
	devices, err := models.GetAllDevices()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		c.HTML(http.StatusOK, "devices.html", gin.H{
			"title":   "Devices",
			"devices": []models.Device{},
			"error":   "Failed to load devices from server",
		})
		return
	}

	c.HTML(http.StatusOK, "devices.html", gin.H{
		"title":   "Devices",
		"devices": devices,
	})
}

// ==================== REST API Endpoints ====================

// DeviceGet - GET /api/devices/:id
func DeviceGet(c *gin.Context) {
	deviceID := c.Param("id")

	device, err := models.GetDeviceByID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"device":  device,
	})
}

// DeviceCreate - POST /api/devices
func DeviceCreate(c *gin.Context) {
	var device models.Device

	// Parse form data
	device.DeviceID = c.PostForm("device_id")
	device.Name = c.PostForm("name")
	device.Status = c.PostForm("status")

	// Parse latitude & longitude
	latStr := c.DefaultPostForm("latitude", "1.1045")
	lngStr := c.DefaultPostForm("longitude", "104.0305")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		log.Printf("Error parsing latitude: %v", err)
		lat = 1.1045 // Default Batam
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		log.Printf("Error parsing longitude: %v", err)
		lng = 104.0305 // Default Batam
	}

	device.Latitude = lat
	device.Longitude = lng

	// Validasi
	if device.Name == "" || device.DeviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Name and Device ID are required",
		})
		return
	}

	_, err = models.CreateDevice(device)
	if err != nil {
		log.Printf("Error creating device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create device: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device added successfully",
		"device":  device,
	})
}

// DeviceUpdate - PUT /api/devices/:id
func DeviceUpdate(c *gin.Context) {
	deviceID := c.Param("id")

	var device models.Device
	device.DeviceID = deviceID
	device.Name = c.PostForm("name")
	device.Status = c.PostForm("status")

	// Parse latitude & longitude
	latStr := c.DefaultPostForm("latitude", "1.1045")
	lngStr := c.DefaultPostForm("longitude", "104.0305")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		log.Printf("Error parsing latitude: %v", err)
		lat = 1.1045
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		log.Printf("Error parsing longitude: %v", err)
		lng = 104.0305
	}

	device.Latitude = lat
	device.Longitude = lng

	// Validasi
	if device.Name == "" || device.DeviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Name and Device ID are required",
		})
		return
	}

	err = models.UpdateDevice(deviceID, device)
	if err != nil {
		log.Printf("Error updating device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update device: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device updated successfully",
	})
}

// DeviceDelete - DELETE /api/devices/:id
func DeviceDelete(c *gin.Context) {
	deviceID := c.Param("id")

	err := models.DeleteDevice(deviceID)
	if err != nil {
		log.Printf("Error deleting device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete device: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device deleted successfully",
	})
}
