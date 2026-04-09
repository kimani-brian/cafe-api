package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // The "-" ensures we NEVER accidentally send the hash in a JSON response
	Role         string `json:"role"`
}

// LoginRequest is the JSON payload we expect from the frontend
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
