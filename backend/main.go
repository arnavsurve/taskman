package main

import (
	"log"
	"net/http"

	"github.com/arnavsurve/taskman/backend/api"
	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	store, err := db.NewPostgresStore()
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
	authRoutes.Use(utils.AuthMiddleware())
	{
		authRoutes.PUT("/user/:id", func(ctx *gin.Context) {
			api.EditUser(ctx, store)
		})
		authRoutes.POST("/table/:id", func(ctx *gin.Context) {
			api.HandleCreateTasksTable(ctx, store)
		})

		authRoutes.POST("/task/:id/:table", func(ctx *gin.Context) {
			api.HandleCreateTask(ctx, store)
		})
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/login", func(ctx *gin.Context) {
		api.Login(ctx, store)
	})
	r.POST("/user", func(ctx *gin.Context) {
		api.AddUser(ctx, store)
	})
	r.GET("/user/:id", func(ctx *gin.Context) {
		api.GetUserByID(ctx, store)
	})

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
