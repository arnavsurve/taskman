package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plaintext password and returns a password hash string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword checks whether a plaintext password matches the requested account's hashed password in the DB
func VerifyPassword(username, password string, store *PostgresStore) (bool, error) {
	var hashedPassword string

	query := `SELECT password from accounts where username = $1`
	err := store.db.QueryRow(query, username).Scan(&hashedPassword)
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
