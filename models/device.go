package models

import "time"

// Device represents a tracking device
type Device struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	DeviceID  string    `json:"device_id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Location  string    `json:"location"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Temporary in-memory storage (ganti dengan database nanti)
var devices = []Device{
	{
		ID:        1,
		Name:      "Tracker Alpha",
		DeviceID:  "TRK-001",
		Type:      "GPS Tracker",
		Status:    "active",
		Location:  "sesuai koordinat nnt",
		Latitude:  1.1182536894443547,
		Longitude: 104.0469369239214,
		LastSeen:  time.Now(),
	},
	{
		ID:        2,
		Name:      "Tracker Beta",
		DeviceID:  "TRK-002",
		Type:      "RFID Tag",
		Status:    "active",
		Location:  "sesuai koordinat nnt",
		Latitude:  1.118015018301657,
		Longitude: 104.04717027611171,
		LastSeen:  time.Now().Add(-2 * time.Hour),
	},
	{
		ID:        3,
		Name:      "Tracker Gamma",
		DeviceID:  "TRK-003",
		Type:      "Bluetooth Beacon",
		Status:    "inactive",
		Location:  "sesuai koordinat nnt",
		Latitude:  1.1181732384995504,
		Longitude: 104.04722660250246,
		LastSeen:  time.Now().Add(-24 * time.Hour),
	},
}

var nextID = 4

// GetAllDevices returns all devices
func GetAllDevices() []Device {
	return devices
}

// GetDeviceByID returns a device by ID
func GetDeviceByID(id int) *Device {
	for i := range devices {
		if devices[i].ID == id {
			return &devices[i]
		}
	}
	return nil
}

// CreateDevice adds a new device
func CreateDevice(device Device) Device {
	device.ID = nextID
	nextID++
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	device.LastSeen = time.Now()

	// Default coordinates untuk Batam jika tidak diset
	if device.Latitude == 0 && device.Longitude == 0 {
		device.Latitude = 1.1045
		device.Longitude = 104.0305
	}

	devices = append(devices, device)
	return device
}

// UpdateDevice updates an existing device
func UpdateDevice(id int, updatedDevice Device) bool {
	for i := range devices {
		if devices[i].ID == id {
			updatedDevice.ID = id
			updatedDevice.CreatedAt = devices[i].CreatedAt
			updatedDevice.UpdatedAt = time.Now()

			// Preserve coordinates if not updated
			if updatedDevice.Latitude == 0 && updatedDevice.Longitude == 0 {
				updatedDevice.Latitude = devices[i].Latitude
				updatedDevice.Longitude = devices[i].Longitude
			}

			devices[i] = updatedDevice
			return true
		}
	}
	return false
}

// DeleteDevice removes a device
func DeleteDevice(id int) bool {
	for i := range devices {
		if devices[i].ID == id {
			devices = append(devices[:i], devices[i+1:]...)
			return true
		}
	}
	return false
}
