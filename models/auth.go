package models

import "time"

// User struct, password is ignored while sending back to client
type User struct {
	Uid        string    `json:"uid"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	IsVerified bool      `json:"isVerified"`
	LastLogin  time.Time `json:"lastLogin"`
	CreatedAt  time.Time `json:"createdAt"`
}

// User sends Register Body to the server while registering
type RegisterBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User sends Login Body to the server while logging in
type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
