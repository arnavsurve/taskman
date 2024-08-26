package main

import (
	"log"
	"net/http"
	// "os"

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

	// Initialize DB with an admin user
	// adminUsername := os.Getenv("DB_ADMIN_USERNAME")
	// adminPassword := os.Getenv("DB_ADMIN_PASSWORD")

	// api.AddUser(, store)

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Apply AuthMiddleware to routes that require authentication
	authRoutes := r.Group("/")
	authRoutes.Use(utils.AuthMiddleware())
	{
		authRoutes.GET("/user/:id", func(ctx *gin.Context) {
			api.GetUserByID(ctx, store)
		})
		authRoutes.PUT("/user/:id", func(ctx *gin.Context) {
			api.EditUser(ctx, store)
		})
		authRoutes.POST("/table/:id", func(ctx *gin.Context) {
			api.HandleCreateTasksTable(ctx, store)
		})

		authRoutes.POST("/task/:id/:workspace", func(ctx *gin.Context) {
			api.HandleCreateTask(ctx, store)
		})
		authRoutes.GET("/task/:id/:workspace", func(ctx *gin.Context) {
			api.HandleGetTasks(ctx, store)
		})
		authRoutes.GET("/task/:id/:workspace/:taskId", func(ctx *gin.Context) {
			api.HandleGetTaskByID(ctx, store)
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

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
