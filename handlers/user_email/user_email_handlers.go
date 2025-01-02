package user_email

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_email_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"log"
	"net/url"
	"strconv"
	"time"
)

// @Summary Get all user emails
// @Description Get all user emails
// @Tags user_email
// @Accept json
// @Produce json
// @Success 200 {array} models.TwUserEmail
// @Router /dbms/v1/user_email [get]
func (h *UserEmailHandler) getUserEmails(c *fiber.Ctx) error {
	// Get user_id from query param
	userId := c.Query("user_id")
	if userId != "" {
		var userEmails []models.TwUserEmail
		if result := h.DB.Where("user_id = ?", userId).Find(&userEmails); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(userEmails)
	}
	var userEmails []models.TwUserEmail
	if result := h.DB.Find(&userEmails); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(userEmails)
}

// @Summary Get user emails by user ID
// @Description Get user emails by user ID
// @Tags user_email
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param status query string false "Status"
// @Success 200 {array} models.TwUserEmail
// @Router /dbms/v1/user_email/user/{user_id} [get]
func (h *UserEmailHandler) getUserEmailByUserId(c *fiber.Ctx) error {
	var userEmails []models.TwUserEmail
	userId := c.Params("user_id")
	status := c.Query("status")
	var query *gorm.DB
	if status == "" {
		query = h.DB.Debug().Where("(status IS NULL OR status = 'linked' OR status = 'pending') AND (user_id = ? OR is_linked_to = ?)", userId, userId)
	} else if status == "pending" {
		query = h.DB.Debug().Where("status = ? AND (user_id = ? OR is_linked_to = ?)", status, userId, userId)
	} else {
		query = h.DB.Debug().Where("(status IS NULL OR status = ?) AND (user_id = ? OR is_linked_to = ?)", status, userId, userId)
	}

	if err := query.Find(&userEmails).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if len(userEmails) == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Emails not found")
	}
	return c.JSON(userEmails)
}

// @Summary Create a new user email
// @Description Create a new user email
// @Tags user_email
// @Accept json
// @Produce json
// @Param email body models.TwUserEmail true "User Email"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email [post]
func (h *UserEmailHandler) createUserEmail(ctx *fiber.Ctx) error {
	userEmail := new(models.TwUserEmail)
	if err := ctx.BodyParser(userEmail); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Create(&userEmail); result.Error != nil {
		var driverErr *mysql.MySQLError
		if errors.As(result.Error, &driverErr) && driverErr.Number == 1062 {
			return ctx.Status(fiber.StatusBadRequest).SendString("email already exists")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	// Lấy thông tin người dùng từ cơ sở dữ liệu dựa trên UserId
	var user models.TwUser
	if err := h.DB.Where("id = ?", userEmail.UserId).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// Gán thông tin user vào userEmail để trả về kèm thông tin user
	userEmail.User = user

	return ctx.JSON(userEmail)
}

// updateUserIdInUserEmail godoc
// @Summary Update user_id in tw_user_email by email
// @Description Update user_id in tw_user_email by email
// @Tags user_email
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Param user_id query string true "User ID"
// @Param status query string true "Status"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email [patch]
func (h *UserEmailHandler) updateStatusInUserEmail(c *fiber.Ctx) error {
	var userEmail models.TwUserEmail
	email := c.Query("email")
	status := c.Query("status")
	if err := h.DB.Where("email = ?", email).First(&userEmail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Email not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	// if status = "", then set null to status, and set null to is_linked_to
	if status == "" {
		userEmail.Status = nil
		userEmail.IsLinkedTo = nil
	} else {
		userEmail.Status = &status
	}
	if status == "pending" {
		expiresAt := time.Now().Add(time.Minute * 10)
		userEmail.ExpiresAt = &expiresAt
	} else {
		userEmail.ExpiresAt = nil
	}
	userEmail.DeletedAt = nil

	if result := h.DB.Save(&userEmail); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(userEmail)
}

// updateUserEmailStatus godoc
// @Summary Update user email status
// @Description Update user email status
// @Tags user_email
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Param status query string true "Status"
// @Param target_user_id query string true "Target User ID"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email/status [patch]
func (h *UserEmailHandler) updateUserEmailStatusAndIsLinkedTo(c *fiber.Ctx) error {
	var userEmail models.TwUserEmail
	email := c.Query("email")
	status := c.Query("status")
	targetUserId := c.Query("target_user_id")
	if err := h.DB.Where("email = ?", email).First(&userEmail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Email not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if status == "" {
		userEmail.Status = nil
		userEmail.IsLinkedTo = nil
	} else {
		userEmail.Status = &status
		targetUserIdInt, err := strconv.Atoi(targetUserId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		userEmail.IsLinkedTo = &targetUserIdInt
	}
	if status == "pending" {
		expiresAt := time.Now().Add(time.Minute * 10)
		userEmail.ExpiresAt = &expiresAt
	} else {
		userEmail.ExpiresAt = nil
	}
	if result := h.DB.Save(&userEmail); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(userEmail)
}

// @Summary Delete user email by ID
// @Description Delete user email by ID
// @Tags user_email
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param email query string true "Email"
// @Param status query string true "status"
// @Success 200 {string} string
// @Router /dbms/v1/user_email [delete]
func (h *UserEmailHandler) deleteUserEmail(ctx *fiber.Ctx) error {
	userId := ctx.Query("user_id")
	email := ctx.Query("email")
	status := ctx.Query("status")
	var userEmail models.TwUserEmail
	if result := h.DB.Where("user_id = ? AND email = ? AND status = ?", userId, email, status).First(&userEmail); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("User email not found")
	}
	if result := h.DB.Delete(&models.TwUserEmail{}, "email = ? AND status = 'pending'", email); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.SendString("Email deleted successfully")
}

// @Summary Get user emails by email
// @Description Get user emails by email
// @Tags user_email
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {array} models.TwUserEmail
// @Router /dbms/v1/user_email/email/{email} [get]
func (h *UserEmailHandler) getUserEmailByEmail(c *fiber.Ctx) error {
	var userEmails models.TwUserEmail
	email := c.Params("email")
	emailFix, err1 := url.QueryUnescape(email)
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err1.Error(),
		})
	}

	if err := h.DB.Where("email = ?", emailFix).Find(&userEmails).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if userEmails.Email == "" {
		return c.Status(fiber.StatusNotFound).SendString("Email not found")
	}

	return c.JSON(userEmails)
}

// @Summary Get exact user email by user ID
// @Description Get exact user email by user ID
// @Tags user_email
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email/user_id/{user_id} [get]
func (h *UserEmailHandler) getExactUserEmailByUserId(c *fiber.Ctx) error {
	var userEmails models.TwUserEmail
	userId := c.Params("user_id")

	if err := h.DB.Where("user_id = ?", userId).Find(&userEmails).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if userEmails.Email == "" {
		return c.Status(fiber.StatusNotFound).SendString("Email not found")
	}

	return c.JSON(userEmails)
}

// getUserEmailToCheckBeforeLink godoc
// @Summary Get user email to check before link
// @Description Get user email to check before link
// @Tags user_email
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email/check [get]
func (h *UserEmailHandler) getUserEmailToCheckBeforeLink(c *fiber.Ctx) error {
	var userEmails models.TwUserEmail
	email := c.Query("email")

	if err := h.DB.Where("email = ? AND is_linked_to is not null AND status is not null", email).First(&userEmails).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Email not found and ok to be linked")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(userEmails)
}

// @Summary Search user email
// @Description Search user email
// @Tags user_email
// @Accept json
// @Produce json
// @Param query path string true "Query"
// @Success 200 {array} user_email_dtos.SearchUserEmailResponse
// @Router /dbms/v1/user_email/search/{query} [get]
func (h *UserEmailHandler) searchUserEmail(c *fiber.Ctx) error {
	query := c.Params("query")
	queryFix, err := url.QueryUnescape(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var userEmailInfo []user_email_dtos.SearchUserEmailResponse

	err = h.DB.
		Table("tw_user_emails").
		Select("tw_user_emails.id, tw_user_emails.email, tw_users.first_name, tw_users.last_name, tw_users.profile_picture").
		Joins("JOIN tw_users ON tw_user_emails.email = tw_users.email").
		Where("tw_user_emails.email LIKE ? OR tw_users.first_name LIKE ? OR tw_users.last_name LIKE ?", "%"+queryFix+"%", "%"+queryFix+"%", "%"+queryFix+"%").
		Where("tw_user_emails.deleted_at IS NULL").
		Where("tw_users.deleted_at IS NULL").
		Scan(&userEmailInfo).Error

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy dữ liệu người dùng",
		})
	}

	return c.JSON(userEmailInfo)
}

func (h *UserEmailHandler) getEmailInProgress(c *fiber.Ctx) error {
	var userInfo []user_email_dtos.UserEmailStatusResponse
	scheduleId := c.Params("scheduleId")

	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schedule ID không hợp lệ",
		})
	}

	// Thực hiện truy vấn SQL với nhiều phép JOIN và loại bỏ email trùng lặp
	err := h.DB.Table("tw_schedule_participants AS sp").
		Select(`
        DISTINCT ue.id AS id,
        ue.email,
        u.first_name,
        u.last_name,
        u.profile_picture,
        wu.status,
        wu.is_verified
    `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id = sp.workspace_user_id").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("sp.schedule_id = ?", scheduleId).
		Where("wu.status = 'pending'").
		Where("wu.is_verified = 0").
		Where("sp.deleted_at IS NULL").
		Where("wu.deleted_at IS NULL").
		Where("ue.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Scan(&userInfo).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(userInfo)

}

// clearExpiredUserEmails godoc
// @Summary Clear expired user emails
// @Description Clear expired user emails
// @Tags user_email
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /dbms/v1/user_email/clear-expired [get]
func (h *UserEmailHandler) clearExpiredUserEmails(c *fiber.Ctx) error {
	if result := h.DB.Model(&models.TwUserEmail{}).
		Where("expires_at <= NOW()").
		Updates(map[string]interface{}{
			"status":       nil,
			"is_linked_to": nil,
			"expires_at":   nil,
		}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.Status(fiber.StatusOK).SendString("Expired user emails cleared successfully")
}

// clearStatusRejectedEmail godoc
// @Summary Clear status rejected email
// @Description Clear status rejected email
// @Tags user_email
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {string} string
// @Router /dbms/v1/user_email/clear-rejected [get]
func (h *UserEmailHandler) clearStatusRejectedEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if result := h.DB.Model(&models.TwUserEmail{}).
		Where("email = ?", email).
		Where("status = 'rejected'").
		Updates(map[string]interface{}{
			"status":       nil,
			"is_linked_to": nil,
			"expires_at":   nil,
		}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.Status(fiber.StatusOK).SendString("Rejected user emails cleared successfully")
}
