package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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
	token, err := generateToken(newAccount.ID)
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
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully updated"})
}

func HandleCreateTasksTable(ctx *gin.Context, store *PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil || intId != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	tableName := ctx.Param("table_name")
	err = store.CreateTasksTable(intId, tableName)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(500, "Failed to create tasks table")
		return
	}

	ctx.JSON(http.StatusOK, "Tasks table successfully created")
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
			return
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}

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
		return nil, fmt.Errorf("invalid token")
	}
}

func generateToken(userID int) (string, error) {
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
