package utils

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginFields struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CompletionStatus string

const (
	Todo       CompletionStatus = "TODO"
	InProgress CompletionStatus = "IN_PROGRESS"
	Done       CompletionStatus = "DONE"
)

type Task struct {
	Description      string           `json:"description"`
	DueDate          time.Time        `json:"due_date"`
	CompletionStatus CompletionStatus `json:"completion"`
	AccountId        int              `json:"account_id"`
}

type Table struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
