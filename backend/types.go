package main

import (
	"strings"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewAccount returns an Account object with a hashed password, case-insensitive email,
// and generates a created at value.
func NewAccount(username, password, email string) *Account {
	hashedPassword, _ := HashPassword(password)
	return &Account{
		Username:  username,
		Password:  hashedPassword,
		Email:     strings.ToLower(email),
		CreatedAt: time.Now().UTC(),
	}
}
