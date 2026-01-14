package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os" // Tambahkan import os

	"github.com/gin-gonic/gin"
)

// Ubah dari const menjadi var
var BackendAPI = "http://localhost:8000"

// Tambahkan fungsi init untuk membaca env
func init() {
	if url := os.Getenv("BACKEND_URL"); url != "" {
		BackendAPI = url
	}
}

// Device struct untuk semua controller
type Device struct {
	DeviceID  string  `json:"device_id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

// LoginRequest dan LoginResponse untuk auth
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error,omitempty"`
}

// Helper: Ambil token dari cookie atau context
func GetAuthToken(c *gin.Context) string {
	// Coba ambil dari cookie dulu
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		return token
	}

	// Kalau tidak ada di cookie, ambil dari context
	if token, exists := c.Get("auth_token"); exists {
		return token.(string)
	}

	return ""
}

// Helper: Kirim request ke backend API dengan JWT
func MakeBackendRequest(method, endpoint string, token string, body interface{}) (*http.Response, error) {
	url := BackendAPI + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	// Kirim JWT di Authorization header
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Execute request
	client := &http.Client{}
	return client.Do(req)
}

// Helper: Parse response devices dari backend
func ParseDevicesResponse(resp *http.Response) ([]Device, error) {
	var devices []Device
	var rawResponse interface{}

	json.NewDecoder(resp.Body).Decode(&rawResponse)

	// Cek format response backend
	switch v := rawResponse.(type) {
	case []interface{}:
		// Jika langsung array
		jsonData, _ := json.Marshal(v)
		json.Unmarshal(jsonData, &devices)
	case map[string]interface{}:
		// Jika ada wrapper "data"
		if data, ok := v["data"]; ok {
			jsonData, _ := json.Marshal(data)
			json.Unmarshal(jsonData, &devices)
		}
	}

	return devices, nil
}
