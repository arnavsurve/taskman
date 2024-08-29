package api

import (
	"net/http"
	"strconv"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func HandleCreateWorkspace(ctx *gin.Context, store *db.PostgresStore) {
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
	// workspaceName := ctx.Param("workspace")

	// Binding JSON body to a Workspace struct
	workspace := shared.Workspace{AccountId: userID}
	if err := ctx.ShouldBindJSON(&workspace); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workspace.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace name cannot be empty"})
		return
	}
	if workspace.AccountId < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	workspaceID, err := store.CreateWorkspace(workspace.Name, workspace.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Successfully created workspace", "workspace_id": workspaceID})
}
