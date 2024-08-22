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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
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
