package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const FastAPIBaseURL = "http://localhost:8000"

type Device struct {
	DeviceID  string  `json:"device_id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

/* JWT */

var jwtToken string

func GetJWTToken() (string, error) {
	if jwtToken != "" {
		return jwtToken, nil
	}

	resp, err := http.Post(FastAPIBaseURL+"/auth/token", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	jwtToken = result.AccessToken
	return jwtToken, nil
}

func authRequest(method, url string, body []byte) (*http.Request, error) {
	token, err := GetJWTToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

/* CRUD */

func GetAllDevices() ([]Device, error) {
	resp, err := http.Get(FastAPIBaseURL + "/barang")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var devices []Device
	err = json.NewDecoder(resp.Body).Decode(&devices)
	return devices, err
}

func GetDeviceByID(id string) (*Device, error) {
	resp, err := http.Get(FastAPIBaseURL + "/barang?device_id=" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list []Device
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("device not found: %s", id)
	}

	return &list[0], nil
}

func CreateDevice(device Device) (*Device, error) {
	payload, _ := json.Marshal(device)

	req, err := authRequest("POST", FastAPIBaseURL+"/barang", payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var created Device
	err = json.NewDecoder(resp.Body).Decode(&created)
	return &created, err
}

func UpdateDevice(id string, device Device) error {
	payload, _ := json.Marshal(map[string]string{
		"name": device.Name,
	})

	req, err := authRequest("PATCH", FastAPIBaseURL+"/barang/"+id, payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed: %s", string(body))
	}

	return nil
}

func DeleteDevice(id string) error {
	req, err := authRequest("DELETE", FastAPIBaseURL+"/barang/"+id, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(body))
	}

	return nil
}
