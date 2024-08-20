package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUser(ctx *gin.Context, store *PostgresStore) {
	body := Account{}

	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad input")
		return
	}

	newAccount := NewAccount(body.Username, body.Password, body.Email) // NewAccount returns a hashed password BTW
	query := `
        INSERT INTO accounts(username, password, email, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	_, err = store.db.Exec(query, newAccount.Username, newAccount.Password, newAccount.Email, newAccount.CreatedAt)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the user")
	} else {
		ctx.JSON(http.StatusOK, "User successfully created")
	}
}
