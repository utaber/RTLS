package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FastAPI base URL - sesuaikan dengan server FastAPI kamu
const FastAPIURL = "http://localhost:8000"

// Device represents a tracking device (match dengan Firebase structure)
type Device struct {
	DeviceID  string    `json:"device_id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// GetAllDevices fetches all devices from FastAPI
func GetAllDevices() ([]Device, error) {
	resp, err := http.Get(FastAPIURL + "/barang")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch devices: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return devices, nil
}

// GetDeviceByID returns a device by device_id from FastAPI
func GetDeviceByID(deviceID string) (*Device, error) {
	url := fmt.Sprintf("%s/barang?device_id=%s", FastAPIURL, deviceID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("device not found")
	}

	return &devices[0], nil
}

// CreateDevice adds a new device via FastAPI
func CreateDevice(device Device) (*Device, error) {
	jsonData, err := json.Marshal(device)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal device: %v", err)
	}

	resp, err := http.Post(
		FastAPIURL+"/barang",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Return the created device
	return &device, nil
}

// UpdateDevice updates an existing device via FastAPI
func UpdateDevice(deviceID string, device Device) error {
	// Set device_id
	device.DeviceID = deviceID

	jsonData, err := json.Marshal(device)
	if err != nil {
		return fmt.Errorf("failed to marshal device: %v", err)
	}

	// FastAPI biasanya menggunakan POST untuk update Firebase Realtime DB
	resp, err := http.Post(
		FastAPIURL+"/barang",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to update device: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteDevice removes a device via FastAPI
func DeleteDevice(deviceID string) error {
	// Untuk Firebase Realtime DB, kita perlu tambah endpoint DELETE di FastAPI
	// Atau set data menjadi null
	client := &http.Client{}
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/barang/%s", FastAPIURL, deviceID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete device: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
