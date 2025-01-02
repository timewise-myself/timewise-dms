package workspace

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWorkspaceHandler(router fiber.Router, db *gorm.DB) {
	workspaceHandler := WorkspaceHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", workspaceHandler.getWorkspaces)
	router.Get("/:workspace_id", workspaceHandler.getWorkspaceById)
	router.Delete("/:workspace_id", workspaceHandler.removeWorkspaceById)
	router.Post("/", workspaceHandler.createWorkspace)
	router.Put("/:workspace_id", workspaceHandler.updateWorkspace)
	router.Get("/user/:user_id", workspaceHandler.getWorkspacesByUserId)
	router.Get("/status/:status", workspaceHandler.getWorkspacesByStatus)
	router.Get("/is_active/:is_active", workspaceHandler.getWorkspacesByIsActive)
	router.Get("/email/:email", workspaceHandler.getWorkspacesByEmail)
	router.Get("/filter/workspace", workspaceHandler.filterWorkspaces)

}
