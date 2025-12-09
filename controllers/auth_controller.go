package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginIndex menampilkan halaman login
func LoginIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// LoginSubmit memproses form login
func LoginSubmit(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validasi input
	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login",
			"error": "Email dan password harus diisi!",
		})
		return
	}

	// TODO: Ganti dengan validasi database
	if email == "admin@gmail.com" && password == "123" {
		// TODO: Implementasi session/JWT untuk auth yang proper
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	// Login gagal
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"title": "Login",
		"error": "Email atau password salah!",
	})
}

// Logout menghapus session dan redirect ke login
func Logout(c *gin.Context) {
	// TODO: Hapus session/token
	c.Redirect(http.StatusFound, "/login")
}
