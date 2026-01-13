package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

/* ================================
   HALAMAN LOGIN
================================ */

func (ac *AuthController) LoginIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

/* ================================
   SUBMIT LOGIN - CALL BACKEND API
================================ */

func (ac *AuthController) LoginSubmit(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validasi input sederhana
	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Email dan password wajib diisi",
			"title": "Login",
		})
		return
	}

	// Kirim request ke backend API
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	reqBody, _ := json.Marshal(loginReq)
	resp, err := http.Post(
		BackendAPI+"/auth/login",
		"application/json",
		bytes.NewBuffer(reqBody),
	)

	if err != nil || resp.StatusCode != http.StatusOK {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Email atau password salah",
			"title": "Login",
		})
		return
	}

	// Parse response JWT dari backend
	var loginResp LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)
	resp.Body.Close()

	if loginResp.AccessToken == "" {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Gagal mendapatkan token dari server",
			"title": "Login",
		})
		return
	}

	// Simpan JWT token dari backend ke cookie HttpOnly
	c.SetCookie(
		"auth_token",             // nama cookie
		loginResp.AccessToken,    // JWT dari backend
		int(time.Hour.Seconds()), // 1 jam
		"/",                      // path
		"",                       // domain
		false,                    // Secure (set true untuk HTTPS)
		true,                     // HttpOnly
	)

	// Simpan email ke cookie biasa (untuk reference)
	c.SetCookie(
		"user_email",
		email,
		int(time.Hour.Seconds()),
		"/",
		"",
		false,
		false, // tidak HttpOnly karena hanya info saja
	)

	c.Redirect(http.StatusFound, "/dashboard")
}

/* ================================
   LOGOUT
================================ */

func (ac *AuthController) Logout(c *gin.Context) {
	// Hapus cookies
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.SetCookie("user_email", "", -1, "/", "", false, false)

	c.Redirect(http.StatusFound, "/login")
}
