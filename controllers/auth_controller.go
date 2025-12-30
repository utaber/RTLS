package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/* ================================
   STRUCT USER DARI FIREBASE RTDB
================================ */

type User struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type UsersResponse struct {
	Admin User `json:"Admin"`
}

const FirebaseDBURL = "https://rtlsrks513-tes-default-rtdb.asia-southeast1.firebasedatabase.app"

/* ================================
   AMBIL DATA USER DARI FIREBASE RTDB
================================ */

func getUserFromFirebase(email string) (*User, error) {
	// Query ke Firebase Realtime Database
	resp, err := http.Get(FirebaseDBURL + "/Users.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Cek status response
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("gagal mengambil data dari Firebase")
	}

	var users UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	// Cek apakah email cocok dengan Admin
	if users.Admin.Email == email {
		return &users.Admin, nil
	}

	return nil, errors.New("user not found")
}

/* ================================
   GENERATE SIMPLE TOKEN
================================ */

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

/* ================================
   VALIDASI LOGIN
================================ */

func validateLogin(email, password string) (string, error) {
	// Ambil user dari Firebase
	user, err := getUserFromFirebase(email)
	if err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	// Cek password (PLAINTEXT - tidak aman untuk production!)
	if user.Password != password {
		return "", errors.New("password salah")
	}

	// Generate token sederhana
	token := generateToken()
	return token, nil
}

/* ================================
   HALAMAN LOGIN
================================ */

func LoginIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

/* ================================
   SUBMIT LOGIN
================================ */

func LoginSubmit(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Email dan password wajib diisi",
			"title": "Login",
		})
		return
	}

	//  VALIDASI LOGIN
	token, err := validateLogin(email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err.Error(),
			"title": "Login",
		})
		return
	}

	//  SIMPAN TOKEN DI COOKIE
	c.SetCookie(
		"firebase_token",
		token,
		int(time.Hour.Seconds()), // 1 jam
		"/",
		"",
		false, // ubah ke true kalau pakai HTTPS
		true,  // HttpOnly
	)

	c.SetCookie(
		"user_email",
		email,
		int(time.Hour.Seconds()),
		"/",
		"",
		false,
		false,
	)

	c.Redirect(http.StatusFound, "/dashboard")
}

/* ================================
   LOGOUT
================================ */

func Logout(c *gin.Context) {
	c.SetCookie("firebase_token", "", -1, "/", "", false, true)
	c.SetCookie("user_email", "", -1, "/", "", false, false)
	c.Redirect(http.StatusFound, "/login")
}
