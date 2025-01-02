package workspace_user

//workspace_user_handler.go
import (
	"errors"
	"github.com/gofiber/fiber/v2"
	workspaceUserDtos "github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

type WorkspaceUserHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func (h *WorkspaceUserHandler) getWorkspaceUsers(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	if result := h.DB.Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUserById(c *fiber.Ctx) error {
	var workspaceUser models.TwWorkspaceUser
	workspaceUserId := c.Params("workspace_user_id")

	if err := h.DB.Where("id = ?", workspaceUserId).First(&workspaceUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("WorkspaceUser not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUser)
}

func (h *WorkspaceUserHandler) removeWorkspaceUserById(c *fiber.Ctx) error {
	workspaceUserId := c.Params("workspace_user_id")
	if result := h.DB.Delete(&models.TwWorkspaceUser{}, workspaceUserId); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *WorkspaceUserHandler) createWorkspaceUser(c *fiber.Ctx) error {
	workspaceUser := new(models.TwWorkspaceUser)
	if err := c.BodyParser(workspaceUser); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Create(workspaceUser); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUser)
}

func (h *WorkspaceUserHandler) updateWorkspaceUser(c *fiber.Ctx) error {
	workspaceUserId := c.Params("workspace_user_id")

	workspaceUser := new(models.TwWorkspaceUser)
	if err := c.BodyParser(workspaceUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var existingUser models.TwWorkspaceUser
	if err := h.DB.First(&existingUser, "id = ?", workspaceUserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Workspace user not found",
			})
		}
		// For any other error, return a 500 status
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	existingUser.UserEmailId = workspaceUser.UserEmailId
	existingUser.WorkspaceId = workspaceUser.WorkspaceId
	existingUser.WorkspaceKey = workspaceUser.WorkspaceKey
	existingUser.Role = workspaceUser.Role
	existingUser.Status = workspaceUser.Status
	existingUser.IsActive = workspaceUser.IsActive
	existingUser.IsVerified = workspaceUser.IsVerified
	existingUser.ExtraData = workspaceUser.ExtraData

	if err := h.DB.Omit("deleted_at").Save(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update workspace user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(existingUser)
}

//func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceId(c *fiber.Ctx) error {
//	var workspaceUsers []models.TwWorkspaceUser
//	workspaceId := c.Params("workspace_id")
//
//	if result := h.DB.Where("workspace_id = ?", workspaceId).Find(&workspaceUsers); result.Error != nil {
//		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
//	}
//
//	return c.JSON(workspaceUsers)
//}

// getWorkspaceUsersByWorkspaceId godoc
// @Summary Get workspace users by workspace ID
// @Description Get workspace users by workspace ID
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} workspaceUserDtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceId(c *fiber.Ctx) error {
	var workspaceUsers []workspaceUserDtos.GetWorkspaceUserListResponse
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.email = tw_users.email").
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		Where("tw_users.deleted_at IS NULL").
		Where("tw_workspace_users.workspace_id = ? and tw_users.is_verified = true and tw_users.is_active = true and tw_workspace_users.status = 'joined' and tw_workspace_users.is_active = true and tw_workspace_users.is_verified=true", workspaceId).
		Scan(&workspaceUsers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUsers)
}

// getWorkspaceUsersByWorkspaceIdForManage godoc
// @Summary Get workspace users by workspace ID for manage
// @Description Get workspace users by workspace ID for manage
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} workspaceUserDtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/manage/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceIdForManage(c *fiber.Ctx) error {
	var workspaceUsers []workspaceUserDtos.GetWorkspaceUserListResponse
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.email = tw_users.email").
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		Where("tw_users.deleted_at IS NULL").
		Where("tw_workspace_users.workspace_id = ? ", workspaceId).
		Scan(&workspaceUsers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByUserId(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	userId := c.Params("user_id")

	if result := h.DB.Where("user_id = ?", userId).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceKey(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	workspaceKey := c.Params("workspace_key")

	if result := h.DB.Where("workspace_key = ?", workspaceKey).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByStatus(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	status := c.Params("status")

	if result := h.DB.Where("status = ?", status).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByIsActive(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	isActive := c.Params("is_active")

	if result := h.DB.Where("is_active = ?", isActive).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

// GET /workspaces/email/{email}
// getWorkspacesByEmail godoc
// @Summary Get workspaces by email
// @Description Get workspaces by email
// @Tags workspace
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} models.TwWorkspace
// @Router /dbms/v1/workspace_user/email/{email}/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) getWorkspaceUserByEmailAndWorkspace(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	email := c.Params("email")

	// Decode the email parameter to handle special characters
	emailFix, err1 := url.QueryUnescape(email)
	if err1 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email",
		})
	}

	// Check required parameters
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}

	var TWorkspaceUser models.TwWorkspaceUser

	// Perform the database query
	err := h.DB.Debug().
		Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key, tw_workspace_users.role, tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified, tw_workspace_users.extra_data").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id = tw_user_emails.id").
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		Where("tw_user_emails.email = ? AND tw_workspace_users.workspace_id = ?", emailFix, workspaceId).
		Scan(&TWorkspaceUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(TWorkspaceUser)
}

// GetWorkspaceUserInvitationList godoc
// @Summary Get workspace user invitation list
// @Description Get workspace user invitation list
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} workspaceUserDtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/invitation/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) GetWorkspaceUserInvitationList(c *fiber.Ctx) error {
	var workspaceUsers []workspaceUserDtos.GetWorkspaceUserListResponse
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.email = tw_users.email").
		Where("tw_workspace_users.workspace_id = ? and tw_users.is_verified = true and tw_users.is_active = false and tw_workspace_users.status != 'joined' and tw_workspace_users.is_active = true ", workspaceId).
		Scan(&workspaceUsers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUsers)
}

// deleteWorkspaceUser godoc
// @Summary Delete workspace user
// @Description Delete workspace user
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param workspace_user_id path string true "Workspace User ID"
// @Success 200 {object} fiber.Map
// @Router /dbms/v1/workspace_user/{workspace_user_id}/workspace/{workspace_id} [delete]
func (h *WorkspaceUserHandler) DeleteWorkspaceUser(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace is required",
		})
	}
	workspaceUserId := c.Params("workspace_user_id")
	if workspaceUserId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace User is required",
		})
	}
	var workspaceUser models.TwWorkspaceUser
	err := h.DB.Where("id = ? and workspace_id = ?", workspaceUserId, workspaceId).First(&workspaceUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if result := h.DB.Model(&workspaceUser).
		Update("deleted_at", gorm.Expr("NOW()")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// UpdateRole godoc
// @Summary Update role of workspace user
// @Description Update role of workspace user
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param workspace_user body workspaceUserDtos.UpdateWorkspaceUserRoleRequest true "Update role request"
// @Success 200 {object} fiber.Map
// @Router /dbms/v1/workspace_user/role/workspace/{workspace_id} [put]
func (h *WorkspaceUserHandler) UpdateRole(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace is required",
		})
	}
	var workspaceUserRequest workspaceUserDtos.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&workspaceUserRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var workspaceUser models.TwWorkspaceUser
	err := h.DB.
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id = tw_user_emails.id").
		Where("tw_user_emails.email = ? and tw_workspace_users.workspace_id = ?", workspaceUserRequest.Email, workspaceId).
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		First(&workspaceUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if workspaceUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workspace User not found",
		})
	}

	if result := h.DB.Model(&workspaceUser).
		Updates(map[string]interface{}{
			"role":       workspaceUserRequest.Role,
			"updated_at": gorm.Expr("NOW()"),
		}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role updated successfully",
	})
}

// verifyMemberInvitationRequest godoc
// @Summary Verify member's request invitation
// @Description Verify member's request invitation
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} fiber.Map
// @Router /dbms/v1/workspace_user/verify-invitation/workspace/{workspace_id}/email/{email} [put]
func (h *WorkspaceUserHandler) VerifyMemberInvitationRequest(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace is required",
		})
	}
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	var workspaceUser models.TwWorkspaceUser
	err := h.DB.Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id = tw_user_emails.id").
		Where("tw_user_emails.email = ? AND tw_workspace_users.workspace_id = ?", email, workspaceId).
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		First(&workspaceUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if workspaceUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workspace User not found",
		})
	}

	if result := h.DB.Model(&workspaceUser).
		Updates(map[string]interface{}{
			"is_verified": true,
			"updated_at":  gorm.Expr("NOW()"),
		}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Invitation verified successfully",
	})
}

// disproveMemberInvitationRequest godoc
// @Summary Disprove member's request invitation
// @Description Disprove member's request invitation
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} fiber.Map
// @Router /dbms/v1/workspace_user/disprove-invitation/workspace/{workspace_id}/email/{email} [put]
func (h *WorkspaceUserHandler) DisproveMemberInvitationRequest(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace is required",
		})
	}
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	var workspaceUser models.TwWorkspaceUser
	err := h.DB.Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id = tw_user_emails.id").
		Where("tw_user_emails.email = ? AND tw_workspace_users.workspace_id = ?", email, workspaceId).
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		First(&workspaceUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if workspaceUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workspace User not found",
		})
	}

	if result := h.DB.Model(&workspaceUser).
		Updates(map[string]interface{}{
			"status":     "removed",
			"updated_at": gorm.Expr("NOW()"),
		}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Invitation disproved successfully",
	})
}

// UpdateWorkspaceUserStatus godoc
// @Summary Update workspace user status
// @Description Update workspace user status
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_user_id path string true "Workspace User ID"
// @Param workspace_user body models.TwWorkspaceUser true "Update status request"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /dbms/v1/workspace_user/update-status/{workspace_user_id} [put]
func (h *WorkspaceUserHandler) UpdateWorkspaceUserStatus(ctx *fiber.Ctx) error {
	workspace_user_id := ctx.Params("workspace_user_id")
	if workspace_user_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_user_id is required",
		})
	}
	var workspaceUser models.TwWorkspaceUser
	if err := h.DB.Where("id = ?", workspace_user_id).First(&workspaceUser).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if workspaceUser.ID == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workspace User not found",
		})
	}
	var workspaceUserRequest models.TwWorkspaceUser
	if err := ctx.BodyParser(&workspaceUserRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Model(&workspaceUser).
		Updates(map[string]interface{}{
			"status":     workspaceUserRequest.Status,
			"updated_at": gorm.Expr("NOW()"),
			"role":       workspaceUserRequest.Role,
		}); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(workspaceUser)

}

// GetWorkspaceUserInfoById godoc
// @Summary Get workspace user info by ID
// @Description Get workspace user info by ID
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_user_id path string true "Workspace User ID"
// @Success 200 {object} workspaceUserDtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/{workspace_user_id}/info [get]
func (h *WorkspaceUserHandler) GetWorkspaceUserInfoById(ctx *fiber.Ctx) error {

	workspace_user_id := ctx.Params("workspace_user_id")
	if workspace_user_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_user_id is required",
		})
	}
	var workspaceUser workspaceUserDtos.GetWorkspaceUserListResponse
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.email = tw_users.email").
		Where("tw_workspace_users.id = ? ", workspace_user_id).
		Scan(&workspaceUser).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.JSON(workspaceUser)
}

// GetWspUserByUserEmailId godoc
// @Summary Get workspace user ID by user email ID
// @Description Get workspace user ID by user email ID
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param user_email_ids body []string true "List of User Email IDs"
// @Success 200 {array} models.TwWorkspaceUser
// @Router /dbms/v1/workspace_user/user_email_id [POST]
func (h *WorkspaceUserHandler) GetWspUserByUserEmailId(c *fiber.Ctx) error {
	var userEmailIds []string
	if err := c.BodyParser(&userEmailIds); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	var workspaces []models.TwWorkspaceUser
	if err := h.DB.Model(&models.TwWorkspaceUser{}).
		Where("user_email_id IN (?)", userEmailIds).
		Scan(&workspaces).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(workspaces)
}

// UpdateWorkspaceUserStatusByEmailAndWorkspace godoc
// @Summary Update workspace user status by email and workspace
// @Description Update workspace user status by email and workspace
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Param status path string true "Status"
// @Param isActive path string true "Is Active"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /dbms/v1/workspace_user/update-status/email/{email}/workspace/{workspace_id}/status/{status}/isActive/{isActive} [put]
func (h *WorkspaceUserHandler) UpdateWorkspaceUserStatusByEmailAndWorkspace(ctx *fiber.Ctx) error {
	email := ctx.Params("email")
	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}
	workspace_id := ctx.Params("workspace_id")
	if workspace_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	status := ctx.Params("status")
	if status == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status is required",
		})
	}
	isActive := ctx.Params("isActive")
	if isActive == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isActive is required",
		})
	}
	isActiveBool, err := strconv.ParseBool(isActive)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isActive must be a boolean",
		})
	}
	isVerified := ctx.Query("is_verified")
	if isVerified == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isVerified is required",
		})
	}
	isVerifiedBool, err := strconv.ParseBool(isVerified)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "isVerified must be a boolean",
		})
	}
	var workspaceUser models.TwWorkspaceUser
	err = h.DB.Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id = tw_user_emails.id").
		Where("tw_user_emails.email = ? AND tw_workspace_users.workspace_id = ?", email, workspace_id).
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		First(&workspaceUser).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if workspaceUser.ID == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workspace User not found",
		})
	}
	if result := h.DB.Model(&workspaceUser).
		Updates(map[string]interface{}{
			"status":      status,
			"updated_at":  gorm.Expr("NOW()"),
			"is_active":   isActiveBool,
			"is_verified": isVerifiedBool,
		}); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(workspaceUser)
}

// GetWorkspaceUserInvitationNotVerifiedList godoc
// @Summary Get workspace user invitation not verified list
// @Description Get workspace user invitation not verified list
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} workspaceUserDtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/invitation_not_verified/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) GetWorkspaceUserInvitationNotVerifiedList(ctx *fiber.Ctx) error {

	var workspaceUsers []workspaceUserDtos.GetWorkspaceUserListNotVerifiedResponse
	workspaceId := ctx.Params("workspace_id")
	if workspaceId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.email = tw_users.email").
		Where("tw_workspace_users.workspace_id = ? and tw_workspace_users.is_verified = false and tw_users.is_verified=true and tw_users.is_active = true", workspaceId).
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		Where("tw_users.deleted_at IS NULL").
		Where("tw_workspace_users.status <> 'declined' AND tw_workspace_users.status <> 'removed'").
		Scan(&workspaceUsers).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.JSON(workspaceUsers)
}

// GetExistingLinkedWorkspaceUser godoc
// @Summary Get existing linked workspace user
// @Description Get existing linked workspace user
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} models.TwWorkspaceUser
// @Router /dbms/v1/workspace_user/check-existing/email/{email}/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) GetExistingLinkedWorkspaceUser(c *fiber.Ctx) error {
	email := c.Params("email")
	decodedEmail, err := url.QueryUnescape(email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}

	var workspaceUsers []string
	query := `
		WITH root_user_id AS (
			SELECT COALESCE(is_linked_to, user_id) AS root_user_id
			FROM tw_user_emails
			WHERE email = ?
		)
		SELECT tue.email as email
		FROM tw_user_emails tue
		JOIN tw_workspace_users twu ON tue.id = twu.user_email_id
		WHERE (tue.user_id = (SELECT root_user_id FROM root_user_id) OR
		       tue.is_linked_to = (SELECT root_user_id FROM root_user_id)) AND twu.workspace_id = ? AND twu.is_active = true AND twu.deleted_at IS NULL
	`

	if err := h.DB.Raw(query, decodedEmail, workspaceID).Scan(&workspaceUsers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(workspaceUsers)
}
