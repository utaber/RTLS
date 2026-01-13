package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	// DB dihapus
}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

/* ================================
   SHOW DEVICES - Render HTML dengan data dari backend
================================ */

func (dc *DeviceController) ShowDevices(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.Redirect(302, "/login")
		return
	}

	// Call backend API untuk get all devices
	resp, err := MakeBackendRequest("GET", "/barang", token, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching devices from backend: %v, Status: %v", err, resp.StatusCode)
		c.HTML(http.StatusOK, "devices.html", gin.H{
			"title":   "Devices",
			"devices": []Device{},
		})
		return
	}

	// Parse response menggunakan helper
	devices, _ := ParseDevicesResponse(resp)
	resp.Body.Close()

	c.HTML(http.StatusOK, "devices.html", gin.H{
		"title":   "Devices",
		"devices": devices,
	})
}

/* ================================
   DEVICE GET - Ambil detail device by ID
================================ */

func (dc *DeviceController) DeviceGet(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	deviceID := c.Param("id")

	// Jika ada :id param, ambil device spesifik
	if deviceID != "" {
		resp, err := MakeBackendRequest("GET", "/barang/"+deviceID, token, nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Device not found",
			})
			return
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		c.JSON(http.StatusOK, result)
		return
	}

	// Jika tidak ada :id param, ambil semua devices
	resp, err := MakeBackendRequest("GET", "/barang", token, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching devices from backend: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"data": []Device{},
		})
		return
	}

	devices, _ := ParseDevicesResponse(resp)
	resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"data": devices})
}

/* ================================
   DEVICE CREATE - REST API endpoint
================================ */

func (dc *DeviceController) DeviceCreate(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Device name is required",
		})
		return
	}

	// Format data sesuai dengan backend
	backendData := map[string]interface{}{
		"name":      input.Name,
		"latitude":  0,
		"longitude": 0,
		"status":    "Terdeteksi",
	}

	// Call backend API - gunakan /barang sesuai backend
	resp, err := MakeBackendRequest("POST", "/barang", token, backendData)
	if err != nil {
		log.Printf("Error creating device in backend: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to connect to backend",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %v", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Backend rejected the request",
		})
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Tambah success flag
	result["success"] = true

	c.JSON(http.StatusCreated, result)
}

/* ================================
   DEVICE UPDATE - REST API endpoint
================================ */

func (dc *DeviceController) DeviceUpdate(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	deviceID := c.Param("id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Device ID is required",
		})
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
		})
		return
	}

	// Format update data
	updateData := map[string]interface{}{
		"name": input.Name,
	}

	// Backend menggunakan PATCH
	resp, err := MakeBackendRequest("PATCH", "/barang/"+deviceID, token, updateData)
	if err != nil {
		log.Printf("Error updating device in backend: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to connect to backend",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %v", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Backend rejected the request",
		})
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Tambah success flag
	result["success"] = true

	c.JSON(http.StatusOK, result)
}

/* ================================
   DEVICE DELETE - REST API endpoint
================================ */

func (dc *DeviceController) DeviceDelete(c *gin.Context) {
	token := GetAuthToken(c)
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	deviceID := c.Param("id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Device ID is required",
		})
		return
	}

	// Call backend API - DELETE /barang/:id
	resp, err := MakeBackendRequest("DELETE", "/barang/"+deviceID, token, nil)
	if err != nil {
		log.Printf("Error deleting device in backend: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to connect to backend",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %v", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Backend rejected the request",
		})
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Tambah success flag
	result["success"] = true

	c.JSON(http.StatusOK, result)
}
