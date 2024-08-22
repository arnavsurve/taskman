package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login authenticates a user with their username and password and returns a JWT for the session
func Login(ctx *gin.Context, store *PostgresStore) {
	credentials := LoginFields{}
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify password
	valid, err := VerifyPassword(credentials.Username, credentials.Password, store)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Retrieve user ID
	var userID int
	query := `select id from accounts where username = $1`
	err = store.db.QueryRow(query, credentials.Username).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "userid": userID})
}

// AddUser adds a user to the database
func AddUser(ctx *gin.Context, store *PostgresStore) {
	account := Account{}
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccount := NewAccount(account.Username, account.Password, account.Email) // NewAccount returns a hashed password BTW
	query := `
        INSERT INTO accounts(username, password, email, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var userID int
	err := store.db.QueryRow(query, newAccount.Username, newAccount.Password, newAccount.Email, newAccount.CreatedAt).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(newAccount.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUserByID retrieves a user's id, username, email, and creation date from the database by ID
func GetUserByID(ctx *gin.Context, store *PostgresStore) (*Account, error) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Invalid ID entered. Enter an integer value")
	}

	row := store.db.QueryRow(`select id,
							username,
							email,
							created_at 
							from accounts where id = $1`, intId)

	account := Account{}
	err = row.Scan(&account.ID, &account.Username, &account.Email, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatusJSON(404, "ID does not exist in accounts")
		} else {
			fmt.Println(err)
			ctx.AbortWithStatusJSON(500, "Internal server error")
		}
		return nil, err
	}
	ctx.JSON(http.StatusOK, account)
	return &account, nil
}

func EditUser(ctx *gin.Context, store *PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil || userID != intId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	account := Account{}
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `update accounts set username = $1, email = $2 where id = $3`
	_, err = store.db.Exec(query, account.Username, account.Email, intId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully updated"})
}

func HandleCreateTasksTable(ctx *gin.Context, store *PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	table := Table{}
	if err := ctx.ShouldBindJSON(&table); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if table.Name == "" {
		ctx.AbortWithStatusJSON(500, "Table name cannot be empty")
		return
	}

	name, err := store.CreateTasksTable(requestedID, table.Name)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(500, "Failed to create table")
		return
	}

	ctx.JSON(http.StatusOK, name)
}
