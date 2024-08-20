package main

import (
	"net/http"

	"github.com/arnavsurve/taskman/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	database.ConnectDatabase()

	router.GET("/ping", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
