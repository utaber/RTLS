package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	DB *db.Client
}

// Device struct untuk response
type Device struct {
	DeviceID  string  `json:"device_id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

func NewDeviceController(dbClient *db.Client) *DeviceController {
	return &DeviceController{DB: dbClient}
}

// ShowDevices menampilkan halaman devices
func (dc *DeviceController) ShowDevices(c *gin.Context) {
	ctx := context.Background()
	ref := dc.DB.NewRef("Barang")

	var devicesMap map[string]map[string]interface{}
	if err := ref.Get(ctx, &devicesMap); err != nil {
		log.Printf("Error fetching devices: %v", err)
		c.HTML(http.StatusOK, "devices.html", gin.H{
			"title":   "Devices",
			"devices": []Device{},
		})
		return
	}

	// Convert map to slice of Device structs
	deviceList := []Device{}
	for deviceID, deviceData := range devicesMap {
		status, _ := deviceData["status"].(string)
		name, _ := deviceData["name"].(string)
		lat, _ := deviceData["latitude"].(float64)
		lng, _ := deviceData["longitude"].(float64)

		deviceList = append(deviceList, Device{
			DeviceID:  deviceID,
			Name:      name,
			Latitude:  lat,
			Longitude: lng,
			Status:    status,
		})
	}

	c.HTML(http.StatusOK, "devices.html", gin.H{
		"title":   "Devices",
		"devices": deviceList,
	})
}

// DeviceGet mengambil detail device berdasarkan ID
func (dc *DeviceController) DeviceGet(c *gin.Context) {
	deviceID := c.Param("id")
	ctx := context.Background()
	ref := dc.DB.NewRef("Barang/" + deviceID)

	var device map[string]interface{}
	if err := ref.Get(ctx, &device); err != nil || device == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"device_id": deviceID,
			"name":      device["name"],
			"latitude":  device["latitude"],
			"longitude": device["longitude"],
			"status":    device["status"],
		},
	})
}

// DeviceCreate membuat device baru
func (dc *DeviceController) DeviceCreate(c *gin.Context) {
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

	ctx := context.Background()

	// Generate device ID
	deviceID := generateDeviceID(dc.DB)

	// Default values untuk device baru
	ref := dc.DB.NewRef("Barang/" + deviceID)
	err := ref.Set(ctx, map[string]interface{}{
		"name":      input.Name,
		"latitude":  1.1045,   // Default Batam
		"longitude": 104.0305, // Default Batam
		"status":    "Tidak_Terdeteksi",
	})

	if err != nil {
		log.Printf("Error creating device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create device",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   true,
		"message":   "Device created successfully",
		"device_id": deviceID,
	})
}

// DeviceUpdate update device berdasarkan ID
func (dc *DeviceController) DeviceUpdate(c *gin.Context) {
	deviceID := c.Param("id")

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

	ctx := context.Background()
	ref := dc.DB.NewRef("Barang/" + deviceID)

	// Cek apakah device exist
	var existing map[string]interface{}
	if err := ref.Get(ctx, &existing); err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	// Update hanya field name
	updates := map[string]interface{}{
		"name": input.Name,
	}

	if err := ref.Update(ctx, updates); err != nil {
		log.Printf("Error updating device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device updated successfully",
	})
}

// DeviceDelete hapus device berdasarkan ID
func (dc *DeviceController) DeviceDelete(c *gin.Context) {
	deviceID := c.Param("id")
	ctx := context.Background()
	ref := dc.DB.NewRef("Barang/" + deviceID)

	// Cek apakah device exist
	var existing map[string]interface{}
	if err := ref.Get(ctx, &existing); err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	if err := ref.Delete(ctx); err != nil {
		log.Printf("Error deleting device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete device",
		})
		return
	}

	// Optional: Tambahkan device ID ke reusable queue
	reusableRef := dc.DB.NewRef("meta/reusable_ids")
	var reusableIDs []string
	if err := reusableRef.Get(ctx, &reusableIDs); err == nil && reusableIDs != nil {
		reusableIDs = append(reusableIDs, deviceID)
		reusableRef.Set(ctx, reusableIDs)
	} else {
		reusableRef.Set(ctx, []string{deviceID})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device deleted successfully",
	})
}

// Helper function untuk generate device ID dengan FIFO reuse
func generateDeviceID(dbClient *db.Client) string {
	ctx := context.Background()
	reusableRef := dbClient.NewRef("meta/reusable_ids")
	counterRef := dbClient.NewRef("meta/device_counter")

	// Cek apakah ada reusable IDs
	var reusableIDs []string
	if err := reusableRef.Get(ctx, &reusableIDs); err == nil && len(reusableIDs) > 0 {
		// Ambil ID pertama (FIFO)
		deviceID := reusableIDs[0]

		// Hapus dari list
		if len(reusableIDs) > 1 {
			reusableRef.Set(ctx, reusableIDs[1:])
		} else {
			reusableRef.Delete(ctx)
		}

		log.Printf("Reusing device ID: %s", deviceID)
		return deviceID
	}

	// Jika tidak ada reusable ID, generate baru
	var counter int
	if err := counterRef.Get(ctx, &counter); err != nil {
		counter = 0
	}

	counter++
	counterRef.Set(ctx, counter)

	deviceID := fmt.Sprintf("BOX-%03d", counter)
	log.Printf("Generated new device ID: %s", deviceID)
	return deviceID
}
