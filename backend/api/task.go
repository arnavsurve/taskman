package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// HandleCreateTask calls the CreateTask function with user ID and target table passed as URL parameters.
// Task attributes are read from Gin context (JSON request body) and passed to CreateTask.
// A table's name is comprised of the user's ID and name of the workspace.
func HandleCreateTask(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Checking if the requested table for the task exists
	workspaceID := ctx.Param("workspaceId")
	exists, err := store.WorkspaceExists(workspaceID)
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace does not exist"})
		return
	}

	// Binding JSON body to a Task struct
	task := shared.Task{}
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if task.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task description cannot be empty"})
		return
	}

	name, err := store.CreateTask(task.Name, workspaceID, task.Description, task.DueDate, task.CompletionStatus, intRequestedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating task"})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Successfully created task " + name})
}

// HandleGetTasks takes url parameters ID and workspace name. It calls GetTasks and returns a
// JSON object holding task JSON objects.
func HandleGetTasks(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Checking if the requested workspace for the task exists
	workspaceId := ctx.Param("workspaceId")
	exists, err := store.WorkspaceExists(workspaceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if workspace exists"})
		fmt.Println(err)
		return
	}
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace does not exist"})
		return
	}

	tasks, err := store.GetTasks(requestedID, workspaceId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// HandleGetTaskByID takes a userID, workspace name, and task ID as URL parameters and returns a task
func HandleGetTaskByID(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify existence of the requested workspace
	workspaceId := ctx.Param("workspaceId")
	exists, err := store.WorkspaceExists(workspaceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if table exists"})
		return
	}
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Table does not exist"})
		return
	}

	taskID := ctx.Param("taskId")
	task, err := store.GetTaskByID(taskID, workspaceId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task does not exist"})
			return
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func HandleUpdateTaskByID(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify existence of the requested workspace
	workspaceId := ctx.Param("workspaceId")
	exists, err := store.WorkspaceExists(workspaceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if workpace exists"})
		return
	}
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace does not exist"})
		return
	}

	taskID := ctx.Param("taskId")

	// Verify existence of requested task
	exists, err = store.TaskExists(workspaceId, taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if task exists"})
		return
	}
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task does not exist"})
		return
	}

	task := shared.Task{}
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.Name == "" || task.Description == "" || task.CompletionStatus == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task fields cannot be empty"})
		return
	}

	err = store.UpdateTaskByID(taskID, workspaceId, task.Name, task.Description, task.DueDate, task.CompletionStatus)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task successfully updated"})
}

// HandleDeleteTaskByID takes a user ID, workspace name, and task ID as URL parameters and
// calls DeleteTaskByID on the target table
func HandleDeleteTaskByID(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify existence of the requested workspace
	workspaceName := ctx.Param("workspace")
	requestedTable := fmt.Sprintf("t_%s_%s", requestedID, workspaceName)

	exists, err := store.TableExists(requestedTable)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if table exists"})
		return
	}
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Table does not exist"})
		return
	}

	taskID := ctx.Param("taskId")
	err = store.DeleteTaskByID(taskID, requestedTable)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task does not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Successfully deleted task " + taskID})
}
