package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/ping", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/user", func(c *gin.Context) {
		AddUser(c, store)
	})
	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		GetUserByID(c, store, id)
	})

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
