package controllers

import (
	"net/http"
	"rtls_rks513/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DevicesIndex menampilkan daftar semua devices
func DevicesIndex(c *gin.Context) {
	devices := models.GetAllDevices()

	// Ambil success message dari query parameter
	successMsg := c.Query("success")

	c.HTML(http.StatusOK, "devices.html", gin.H{
		"title":   "Devices",
		"devices": devices,
		"success": successMsg,
	})
}

// DeviceGet API endpoint untuk mendapatkan data device (untuk modal)
func DeviceGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid device ID"})
		return
	}

	device := models.GetDeviceByID(id)
	if device != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"device":  device,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
	}
}

// DeviceAddForm menampilkan form tambah device
func DeviceAddForm(c *gin.Context) {
	c.HTML(http.StatusOK, "add_device.html", gin.H{
		"title": "Add Device",
	})
}

// DeviceAdd memproses form tambah device
func DeviceAdd(c *gin.Context) {
	device := models.Device{
		Name:     c.PostForm("name"),
		DeviceID: c.PostForm("device_id"),
		Type:     c.PostForm("type"),
		Status:   c.PostForm("status"),
		Location: c.PostForm("location"),
	}

	models.CreateDevice(device)

	// Redirect dengan success message
	c.Redirect(http.StatusFound, "/devices?success=Device added successfully")
}

// DeviceEditForm menampilkan form edit device
func DeviceEditForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/devices")
		return
	}

	device := models.GetDeviceByID(id)
	if device == nil {
		c.Redirect(http.StatusFound, "/devices")
		return
	}

	c.HTML(http.StatusOK, "edit_device.html", gin.H{
		"title":  "Edit Device",
		"device": device,
	})
}

// DeviceEdit memproses form edit device
func DeviceEdit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/devices")
		return
	}

	device := models.Device{
		Name:     c.PostForm("name"),
		DeviceID: c.PostForm("device_id"),
		Type:     c.PostForm("type"),
		Status:   c.PostForm("status"),
		Location: c.PostForm("location"),
	}

	if models.UpdateDevice(id, device) {
		c.Redirect(http.StatusFound, "/devices?success=Device updated successfully")
	} else {
		c.Redirect(http.StatusFound, "/devices")
	}
}

// DeviceDeleteForm menampilkan konfirmasi hapus device
func DeviceDeleteForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/devices")
		return
	}

	device := models.GetDeviceByID(id)
	if device == nil {
		c.Redirect(http.StatusFound, "/devices")
		return
	}

	c.HTML(http.StatusOK, "delete_device.html", gin.H{
		"title":  "Delete Device",
		"device": device,
	})
}

// DeviceDelete memproses hapus device
func DeviceDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/devices")
		return
	}

	if models.DeleteDevice(id) {
		c.Redirect(http.StatusFound, "/devices?success=Device deleted successfully")
	} else {
		c.Redirect(http.StatusFound, "/devices")
	}
}
