package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// HashPassword takes a plaintext password and returns a password hash string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword checks whether a plaintext password matches the requested account's hashed password in the DB
func VerifyPassword(username, password string, store *db.PostgresStore) (bool, error) {
	var hashedPassword string

	query := `SELECT password from accounts where username = $1`
	err := store.DB.QueryRow(query, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// username not found
			return false, nil
		}
		return false, err
	}

	// compare provided password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

// NewAccount returns an Account object with a hashed password, case-insensitive email,
// and generates a created at value.
func NewAccount(username, password, email string) *shared.Account {
	hashedPassword, _ := HashPassword(password)
	return &shared.Account{
		Username:  username,
		Password:  hashedPassword,
		Email:     strings.ToLower(email),
		CreatedAt: time.Now().UTC(),
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := parseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			fmt.Println(err)
			return
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}

// parseToken takes a JWT token string as input and returns its claims if the token is valid.
func parseToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Invalid token")
	}
}

// GenerateToken generates a new JWT given a user ID
func GenerateToken(userID int) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
