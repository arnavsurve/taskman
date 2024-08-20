package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUser(ctx *gin.Context, store *PostgresStore) {
	body := User{}

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

	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		fmt.Printf("Unable to hash password: %s\n", err)
		return
	}

	fmt.Println(body.Username, body.Password, hashedPassword)

	_, err = store.db.Exec(`
	       INSERT INTO accounts(username, password)
	       values ($1, $2)
	       `, body.Username, hashedPassword)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the user")
	} else {
		ctx.JSON(http.StatusOK, "User successfully created")
	}
}
