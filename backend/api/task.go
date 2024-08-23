package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// HandleCreateTasksTable calls the CreateTasksTable function in db given a user's ID and name of the table.
// On success, returns 400 and created table name.
func HandleCreateTasksTable(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Binding JSON body to a Table struct
	table := shared.Table{}
	if err := ctx.ShouldBindJSON(&table); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if table.Name == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Table name cannot be empty"})
		return
	}

	// Check if table exists
	tableName := fmt.Sprintf("t_%s_%s", requestedID, table.Name)
	if yes, err := store.TableExists(tableName); yes == true {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Table already exists with that name"})
		return
	}

	name, err := store.CreateTasksTable(requestedID, table.Name)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tasks table"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"table_name": name})
}

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
	workspaceName := ctx.Param("workspace")
	requestedTable := fmt.Sprintf("t_%s_%s", requestedID, workspaceName)
	exists, err := store.TableExists(requestedTable)
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Table does not exist"})
		return
	}

	// Binding JSON body to a Task struct
	newTask := shared.Task{}
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newTask.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task description cannot be empty"})
		return
	}

	name, err := store.CreateTask(requestedTable, newTask.Name, newTask.Description, newTask.DueDate, newTask.CompletionStatus, intRequestedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating task"})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"table": name})
}

// HandleGetTasks ... TODO!!!
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

	// Checking if the requested table for the task exists
	workspaceName := ctx.Param("workspace")
	requestedTable := fmt.Sprintf("t_%s_%s", requestedID, workspaceName)
	exists, err := store.TableExists(requestedTable)
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Table does not exist"})
		return
	}

}
