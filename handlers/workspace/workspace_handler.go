package workspace

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
)

type WorkspaceHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

// GetWorkspaces godoc
// @Summary Get all workspaces
// @Description Get all workspaces
// @Tags workspace
// @Accept json
// @Produce json
// @Success 200 {array} models.TwWorkspace
// @Router /dbms/v1/workspace [get]
func (handler *WorkspaceHandler) getWorkspaces(c *fiber.Ctx) error {
	var workspaces []models.TwWorkspace
	if result := handler.DB.Find(&workspaces); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaces)
}

// GET /workspaces/{workspace_id}
// getWorkspaceById godoc
// @Summary Get workspace by ID
// @Description Get workspace by ID
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace_id path int true "Workspace ID"
// @Success 200 {object} models.TwWorkspace
// @Router /dbms/v1/workspace/{workspace_id} [get]
func (handler *WorkspaceHandler) getWorkspaceById(c *fiber.Ctx) error {
	var workspace models.TwWorkspace
	workspaceId := c.Params("workspace_id")

	if err := handler.DB.Where("id = ? and deleted_at IS NULL", workspaceId).First(&workspace).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Workspace not found")
	}
	return c.JSON(workspace)
}

// DELETE /workspaces/{workspace_id}
// removeWorkspaceById godoc
// @Summary Remove workspace by ID
// @Description Remove workspace by ID
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace_id path int true "Workspace ID"
// @Success 204
// @Router /dbms/v1/workspace/{workspace_id} [delete]
func (handler *WorkspaceHandler) removeWorkspaceById(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace ID is required",
		})
	}

	var workspace models.TwWorkspace
	if err := handler.DB.Where("id = ?", workspaceId).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Workspace not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if result := handler.DB.Model(&workspace).Updates(map[string]interface{}{
		"deleted_at": gorm.Expr("NOW()"),
		"is_deleted": true,
	}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// POST /workspaces
// createWorkspace godoc
// @Summary Create workspace
// @Description Create workspace
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace body models.TwWorkspace true "Workspace"
// @Success 200 {object} models.TwWorkspace
// @Router /dbms/v1/workspace [post]
func (handler *WorkspaceHandler) createWorkspace(c *fiber.Ctx) error {
	workspace := new(models.TwWorkspace)
	if err := c.BodyParser(workspace); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := handler.DB.Create(workspace); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspace)
}

// PUT /workspaces
// updateWorkspace godoc
// @Summary Update workspace
// @Description Update workspace
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace body models.TwWorkspace true "Workspace"
// @Success 200 {object} models.TwWorkspace
// @Router /dbms/v1/workspace [put]
func (handler *WorkspaceHandler) updateWorkspace(c *fiber.Ctx) error {
	workspace := new(models.TwWorkspace)
	if err := c.BodyParser(workspace); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Set the UpdatedAt field to the current timestamp using gorm.Expr("NOW()")
	if result := handler.DB.Model(workspace).Update("updated_at", gorm.Expr("NOW()")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	if result := handler.DB.Omit("deleted_at", "created_at").Save(workspace); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspace)
}

// GET /workspaces/user/{user_id}
// getWorkspacesByUserId godoc
// @Summary Get workspaces by user ID
// @Description Get workspaces by user ID
// @Tags workspace
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} []models.TwWorkspace
// @Router /dbms/v1/workspace/user/{user_id} [get]
func (handler *WorkspaceHandler) getWorkspacesByUserId(c *fiber.Ctx) error {
	userId := c.Params("user_id")

	var workspaces []models.TwWorkspace

	// Thực hiện JOIN giữa các bảng để lấy các workspace liên quan đến email
	err := handler.DB.
		Table("tw_workspaces").
		Select("tw_workspaces.id, tw_workspaces.created_at, tw_workspaces.updated_at, tw_workspaces.deleted_at, tw_workspaces.title, tw_workspaces.extra_data, tw_workspaces.description, tw_workspaces.key, tw_workspaces.type, tw_workspaces.is_deleted").
		Joins("JOIN tw_workspace_users ON tw_workspaces.id = tw_workspace_users.workspace_id").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.user_id = tw_users.id").
		Where("tw_users.id = ? and tw_workspace_users.is_active = true and tw_workspace_users.is_verified = true and tw_workspace_users.role != 'guest' and tw_workspace_users.status ='joined'", userId).
		Where("tw_workspaces.deleted_at IS NULL").
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		Where("tw_users.deleted_at IS NULL").
		Where("tw_workspaces.is_deleted = false").
		Scan(&workspaces).Error

	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy workspace",
		})
	}

	return c.Status(http.StatusOK).JSON(workspaces)

}

// GET /workspaces/status/{status}
// getWorkspacesByStatus godoc
// @Summary Get workspaces by status
// @Description Get workspaces by status
// @Tags workspace
// @Accept json
// @Produce json
// @Param status path string true "Status"
// @Success 200 {object} []models.TwWorkspace
// @Router /dbms/v1/workspace/status/{status} [get]
func (handler *WorkspaceHandler) getWorkspacesByStatus(c *fiber.Ctx) error {
	var workspaces []models.TwWorkspace
	status := c.Params("status")
	if result := handler.DB.Where("status = ?", status).Find(&workspaces); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(workspaces)
}

// GET /workspaces/is_active/{is_active}
// getWorkspacesByIsActive godoc
// @Summary Get workspaces by is_active
// @Description Get workspaces by is_active
// @Tags workspace
// @Accept json
// @Produce json
// @Param is_active path string true "Is Active"
// @Success 200 {object} []models.TwWorkspace
// @Router /dbms/v1/workspace/is_active/{is_active} [get]
func (handler *WorkspaceHandler) getWorkspacesByIsActive(c *fiber.Ctx) error {
	var workspaces []models.TwWorkspace
	isActive := c.Params("is_active")
	if result := handler.DB.Where("is_active = ?", isActive).Find(&workspaces); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(workspaces)
}

// getWorkspacesByEmail godoc
// @Summary Get workspaces by email
// @Description Get workspaces by email
// @Tags workspace
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} []models.TwWorkspace
// @Router /dbms/v1/workspace/email/{email} [get]s
func (handler *WorkspaceHandler) getWorkspacesByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	emails, err1 := url.QueryUnescape(email)
	if err1 != nil {
		log.Println("Lỗi khi giải mã email:", err1)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email không hợp lệ",
		})
	}
	var workspaces []models.TwWorkspace

	// Thực hiện JOIN giữa các bảng để lấy các workspace liên quan đến email
	err := handler.DB.
		Table("tw_workspaces").
		Select("tw_workspaces.id, tw_workspaces.created_at, tw_workspaces.updated_at, tw_workspaces.deleted_at, tw_workspaces.title, tw_workspaces.extra_data, tw_workspaces.description, tw_workspaces.key, tw_workspaces.type, tw_workspaces.is_deleted").
		Joins("JOIN tw_workspace_users ON tw_workspaces.id = tw_workspace_users.workspace_id").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Where("tw_user_emails.email = ? and tw_workspace_users.is_active = true and tw_workspace_users.is_verified = true and tw_workspace_users.status ='joined'", emails).
		Where("tw_workspaces.deleted_at IS NULL").
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.deleted_at IS NULL").
		Scan(&workspaces).Error

	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy workspace",
		})
	}

	return c.Status(http.StatusOK).JSON(workspaces)

}

// filterWorkspaces godoc
// @Summary Filter workspaces
// @Description Filter workspaces
// @Tags workspace
// @Accept json
// @Produce json
// @Param email query string false "Email"
// @Param role query string false "Role"
// @Param search query string false "Search"
// @Param sortBy query string false "Sort by"
// @Param order query string false "Order"
// @Param userid query string true "User ID"
// @Success 200 {object} []models.TwWorkspace
// @Router /dbms/v1/workspace/filter/workspace [get]
func (handler *WorkspaceHandler) filterWorkspaces(c *fiber.Ctx) error {
	var workspaces []models.TwWorkspace
	query := handler.DB
	if c.Query("userid") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	query = query.Joins("JOIN tw_workspace_users ON tw_workspaces.id = tw_workspace_users.workspace_id").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id = tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.user_id = tw_users.id").
		Where("tw_users.deleted_at IS NULL").
		Where("tw_workspaces.deleted_at IS NULL").
		Where("tw_workspace_users.is_active = true").
		Where("tw_workspace_users.is_verified = true").
		Where("tw_workspace_users.status = 'joined'").
		Where("tw_workspace_users.role != 'guest'").
		Where("tw_workspace_users.role != 'Guest'").
		Where("tw_workspace_users.deleted_at IS NULL").
		Where("tw_user_emails.user_id = ? or (tw_user_emails.is_linked_to = ? and tw_user_emails.status = 'linked') ", c.Query("userid"), c.Query("userid"))
	// Filter by email
	if email := c.Query("email"); email != "" {
		query = query.
			Where("tw_user_emails.email = ?", email)
	}

	// Filter by role
	if role := c.Query("role"); role != "" {
		query = query.
			Where("tw_workspace_users.role = ?", role)
	}

	// Search by keyword
	if search := c.Query("search"); search != "" {
		query = query.Where("tw_workspaces.title LIKE ? ", "%"+search+"%")
	}

	// Sort by field
	if sortBy := c.Query("sortBy"); sortBy != "" {
		order := c.Query("order", "asc")
		query = query.Order(sortBy + " " + order)
	}

	// Execute the query
	if err := query.Find(&workspaces).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(workspaces)
}
