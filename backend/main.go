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
	r := gin.Default()

	// Apply AuthMiddleware to routes that require authentication
	authRoutes := r.Group("/")
	authRoutes.Use(AuthMiddleware())
	{
		authRoutes.PUT("/user/:id", func(ctx *gin.Context) {
			EditUser(ctx, store)
		})
		authRoutes.POST("/table/:id", func(ctx *gin.Context) {
			HandleCreateTasksTable(ctx, store)
		})
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/login", func(ctx *gin.Context) {
		Login(ctx, store)
	})
	r.POST("/user", func(ctx *gin.Context) {
		AddUser(ctx, store)
	})
	r.GET("/user/:id", func(ctx *gin.Context) {
		GetUserByID(ctx, store)
	})

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
