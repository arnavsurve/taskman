package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func HandleCreateTasksTable(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	table := utils.Table{}
	if err := ctx.ShouldBindJSON(&table); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if table.Name == "" {
		ctx.AbortWithStatusJSON(500, "Table name cannot be empty")
		return
	}

	name, err := store.CreateTasksTable(requestedID, table.Name)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(500, "Failed to create table")
		return
	}

	ctx.JSON(http.StatusOK, name)
}

func HandleCreateTask(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	tableDesc := ctx.Param("table")
	requestedTable := fmt.Sprintf("t_%s_%s", requestedID, tableDesc)
	exists, err := db.TableExists(store, requestedTable)
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Table does not exist"})
		return
	}

	newTask := utils.Task{}
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newTask.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task description cannot be empty"})
		return
	}

	// TODO: there is no way due date is implemented properly bruh fix that shit
	name, err := store.CreateTask(requestedTable, newTask.Description, newTask.DueDate.Format("08-23-2024 08:12 PM"), newTask.CompletionStatus, newTask.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating task"})
	}

	ctx.JSON(http.StatusOK, name)
}
