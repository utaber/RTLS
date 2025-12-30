package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	DB *db.Client
}

func NewDashboardController(dbClient *db.Client) *DashboardController {
	return &DashboardController{DB: dbClient}
}

// ShowDashboard menampilkan halaman dashboard
func (dc *DashboardController) ShowDashboard(c *gin.Context) {
	ctx := context.Background()
	ref := dc.DB.NewRef("Barang")

	var devices map[string]map[string]interface{}
	if err := ref.Get(ctx, &devices); err != nil {
		log.Printf("Error fetching devices: %v", err)
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

	deviceLocations := []gin.H{}
	totalDevices := len(devices)
	activeDevices := 0
	alertsCount := 0

	for deviceID, deviceData := range devices {
		status, _ := deviceData["status"].(string)
		name, _ := deviceData["name"].(string)
		lat, _ := deviceData["latitude"].(float64)
		lng, _ := deviceData["longitude"].(float64)

		if status == "Terdeteksi" {
			activeDevices++
		} else {
			alertsCount++
		}

		deviceLocations = append(deviceLocations, gin.H{
			"deviceID": deviceID,
			"name":     name,
			"lat":      lat,
			"lng":      lng,
			"status":   status,
		})
	}

	activePercentage := 0
	if totalDevices > 0 {
		activePercentage = (activeDevices * 100) / totalDevices
	}

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

// StreamDevices endpoint untuk real-time updates menggunakan SSE
func (dc *DashboardController) StreamDevices(c *gin.Context) {
	ctx := c.Request.Context()

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

	ref := dc.DB.NewRef("Barang")

	// Channel untuk menerima updates
	updateChan := make(chan map[string]interface{}, 10)
	errorChan := make(chan error, 1)

	// Goroutine untuk polling Firebase
	go func() {
		// Kirim data awal
		var devices map[string]map[string]interface{}
		if err := ref.Get(context.Background(), &devices); err != nil {
			errorChan <- err
			return
		}

		updateChan <- map[string]interface{}{"devices": devices}

		// Polling setiap 2 detik untuk perubahan
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		var lastData map[string]map[string]interface{}
		lastData = devices

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var currentDevices map[string]map[string]interface{}
				if err := ref.Get(context.Background(), &currentDevices); err != nil {
					log.Printf("Error polling Firebase: %v", err)
					continue
				}

				// Cek apakah data berubah
				currentJSON, _ := json.Marshal(currentDevices)
				lastJSON, _ := json.Marshal(lastData)

				if string(currentJSON) != string(lastJSON) {
					updateChan <- map[string]interface{}{"devices": currentDevices}
					lastData = currentDevices
				}
			}
		}
	}()

	// Kirim updates ke client
	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected from SSE")
			return
		case err := <-errorChan:
			log.Printf("SSE Error: %v", err)
			fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		case update := <-updateChan:
			devices, ok := update["devices"].(map[string]map[string]interface{})
			if !ok {
				continue
			}

			deviceLocations := []gin.H{}
			for deviceID, deviceData := range devices {
				status, _ := deviceData["status"].(string)
				name, _ := deviceData["name"].(string)
				lat, _ := deviceData["latitude"].(float64)
				lng, _ := deviceData["longitude"].(float64)

				deviceLocations = append(deviceLocations, gin.H{
					"deviceID": deviceID,
					"name":     name,
					"lat":      lat,
					"lng":      lng,
					"status":   status,
				})
			}

			jsonData, err := json.Marshal(deviceLocations)
			if err != nil {
				log.Printf("Error marshaling data: %v", err)
				continue
			}

			// Kirim SSE event
			fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)
			flusher.Flush()
		}
	}
}
