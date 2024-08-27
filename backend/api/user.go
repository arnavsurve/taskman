package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/arnavsurve/taskman/backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login authenticates a user with their username and password and returns a JWT for the session
func Login(ctx *gin.Context, store *db.PostgresStore) {
	credentials := shared.LoginFields{}
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify password
	valid, err := utils.VerifyPassword(credentials.Username, credentials.Password, store)
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
	err = store.DB.QueryRow(query, credentials.Username).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "userid": userID})
}

// AddUser adds a user to the database
func AddUser(ctx *gin.Context, store *db.PostgresStore) {
	account := shared.Account{}
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccount := utils.NewAccount(account.Username, account.Password, account.Email) // NewAccount returns a hashed password BTW
	query := `
        INSERT INTO accounts(username, password, email, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var userID int
	err := store.DB.QueryRow(query, newAccount.Username, newAccount.Password, newAccount.Email, newAccount.CreatedAt).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(newAccount.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUserByID retrieves a user's id, username, email, and creation date from the database by ID
func GetUserByID(ctx *gin.Context, store *db.PostgresStore) (*shared.Account, error) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil || userID != intId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil, err
	}

	row := store.DB.QueryRow(`select id,
							username,
							email,
							created_at
							from accounts where id = $1`, intId)

	account := shared.Account{}
	err = row.Scan(&account.ID, &account.Username, &account.Email, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatusJSON(404, "ID does not exist in accounts")
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return nil, err
	}
	ctx.JSON(http.StatusOK, account)
	return &account, nil
}

func EditUser(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil || userID != intId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	account := shared.Account{}
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `update accounts set username = $1, email = $2 where id = $3`
	_, err = store.DB.Exec(query, account.Username, account.Email, intId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully updated"})
}
