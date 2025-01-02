package workspace_user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWorkspaceUserHandler(router fiber.Router, db *gorm.DB) {
	workspaceUserHandler := WorkspaceUserHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", workspaceUserHandler.getWorkspaceUsers)
	router.Get("/:workspace_user_id", workspaceUserHandler.getWorkspaceUserById)
	router.Delete("/:workspace_user_id", workspaceUserHandler.removeWorkspaceUserById)
	router.Post("/", workspaceUserHandler.createWorkspaceUser)
	router.Put("/:workspace_user_id", workspaceUserHandler.updateWorkspaceUser)
	router.Get("/workspace/:workspace_id", workspaceUserHandler.getWorkspaceUsersByWorkspaceId)
	router.Get("/manage/workspace/:workspace_id", workspaceUserHandler.getWorkspaceUsersByWorkspaceIdForManage)
	router.Get("/user/:user_id", workspaceUserHandler.getWorkspaceUsersByUserId)
	router.Get("/workspace_key/:workspace_key", workspaceUserHandler.getWorkspaceUsersByWorkspaceKey)
	router.Get("/status/:status", workspaceUserHandler.getWorkspaceUsersByStatus)
	router.Get("/is_active/:is_active", workspaceUserHandler.getWorkspaceUsersByIsActive)
	router.Get("/email/:email/workspace/:workspace_id", workspaceUserHandler.getWorkspaceUserByEmailAndWorkspace)
	router.Get("/invitation/workspace/:workspace_id", workspaceUserHandler.GetWorkspaceUserInvitationList)
	router.Delete("/:workspace_user_id/workspace/:workspace_id/", workspaceUserHandler.DeleteWorkspaceUser)
	router.Put("/role/workspace/:workspace_id", workspaceUserHandler.UpdateRole)
	router.Put("/verify-invitation/workspace/:workspace_id/email/:email", workspaceUserHandler.VerifyMemberInvitationRequest)
	router.Put("/disprove-invitation/workspace/:workspace_id/email/:email", workspaceUserHandler.DisproveMemberInvitationRequest)
	router.Put("/update-status/:workspace_user_id", workspaceUserHandler.UpdateWorkspaceUserStatus)
	router.Get("/:workspace_user_id/info", workspaceUserHandler.GetWorkspaceUserInfoById)
	router.Put("/update-status/email/:email/workspace/:workspace_id/status/:status/is_active/:isActive", workspaceUserHandler.UpdateWorkspaceUserStatusByEmailAndWorkspace)
	router.Get("/invitation_not_verified/workspace/:workspace_id", workspaceUserHandler.GetWorkspaceUserInvitationNotVerifiedList)
	router.Get("/check-existing/email/:email/workspace/:workspace_id", workspaceUserHandler.GetExistingLinkedWorkspaceUser)
	router.Post("/user_email_id", workspaceUserHandler.GetWspUserByUserEmailId)
}
