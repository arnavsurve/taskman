package main

import (
	"log"
	"net/http"

	"github.com/arnavsurve/taskman/backend/api"
	"github.com/arnavsurve/taskman/backend/auth"
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

	// initializing GitHub OAuth config
	auth.GithubConfig()

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Apply AuthMiddleware to routes that require authentication
	authRoutes := r.Group("/")
	authRoutes.Use(utils.AuthMiddleware())
	{
		// User routes
		authRoutes.GET("/user/:id", func(ctx *gin.Context) {
			api.GetUserByID(ctx, store)
		})
		authRoutes.PUT("/user/:id", func(ctx *gin.Context) {
			api.EditUser(ctx, store)
		})
		authRoutes.DELETE("/user/:id", func(ctx *gin.Context) {
			api.DeleteUser(ctx, store)
		})

		// Workspace routes
		authRoutes.POST("/workspace/:id", func(ctx *gin.Context) {
			api.HandleCreateWorkspace(ctx, store)
		})
		authRoutes.GET("/workspace/:id", func(ctx *gin.Context) {
			api.HandleListWorkspaces(ctx, store)
		})
		authRoutes.GET("/workspace/:id/:workspaceId", func(ctx *gin.Context) {
			api.HandleGetTasksInWorkspace(ctx, store)
		})
		authRoutes.PUT("/workspace/:id/:workspaceId", func(ctx *gin.Context) {
			api.HandleUpdateWorkspaceByID(ctx, store)
		})
		authRoutes.DELETE("/workspace/:id/:workspaceId", func(ctx *gin.Context) {
			api.HandleDeleteWorkspaceByID(ctx, store)
		})

		// Task routes
		authRoutes.POST("/task/:id/:workspaceId", func(ctx *gin.Context) {
			api.HandleCreateTask(ctx, store)
		})
		authRoutes.GET("/task/:id/:workspaceId/:taskId", func(ctx *gin.Context) {
			api.HandleGetTaskByID(ctx, store)
		})
		authRoutes.PUT("/task/:id/:workspaceId/:taskId", func(ctx *gin.Context) {
			api.HandleUpdateTaskByID(ctx, store)
		})
		authRoutes.DELETE("/task/:id/:workspaceId/:taskId", func(ctx *gin.Context) {
			api.HandleDeleteTaskByID(ctx, store)
		})
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Login routes
	r.GET("/login", func(ctx *gin.Context) {
		api.Login(ctx, store)
	})
	r.POST("/user", func(ctx *gin.Context) {
		api.AddUser(ctx, store)
	})

	r.GET("/oauth2/github", func(ctx *gin.Context) {
		auth.GithubLogin(ctx)
	})
	r.GET("/oauth2/callback", func(ctx *gin.Context) {
		auth.GithubCallback(ctx, store)
	})

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
