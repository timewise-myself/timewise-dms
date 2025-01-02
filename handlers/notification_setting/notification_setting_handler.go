package notification_setting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type NotificationSettingHandler struct {
	DB *gorm.DB
}

// getNotificationSettingByUserId godoc
// @Summary Get notification setting by user id
// @Description Get notification setting by user id
// @Tags notification_setting
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.TwNotificationSettings
// @Router /dbms/v1/notification_setting/{user_id} [get]
func (h NotificationSettingHandler) GetNotificationSettingByUserId(ctx *fiber.Ctx) error {
	id := ctx.Params("user_id")
	var notificationSetting models.TwNotificationSettings
	if result := h.DB.Where("user_id = ?", id).First(&notificationSetting); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	return ctx.JSON(notificationSetting)

}

// CreateNotificationSetting godoc
// @Summary Create notification setting
// @Description Create notification setting
// @Tags notification_setting
// @Accept json
// @Produce json
// @Param notification_setting body models.TwNotificationSettings true "Notification Setting"
// @Success 200 {object} models.TwNotificationSettings
// @Router /dbms/v1/notification_setting [post]
func (h NotificationSettingHandler) CreateNotificationSetting(ctx *fiber.Ctx) error {
	var notificationSetting models.TwNotificationSettings
	if err := ctx.BodyParser(&notificationSetting); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&notificationSetting); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(notificationSetting)

}

// UpdateNotificationSetting godoc
// @Summary Update notification setting
// @Description Update notification setting
// @Tags notification_setting
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param notification_setting body models.TwNotificationSettings true "Notification Setting"
// @Success 200 {object} models.TwNotificationSettings
// @Router /dbms/v1/notification_setting/{user_id} [put]
func (h NotificationSettingHandler) UpdateNotificationSetting(ctx *fiber.Ctx) error {
	id := ctx.Params("user_id")
	var notificationSetting models.TwNotificationSettings
	if result := h.DB.Where("user_id = ?", id).First(&notificationSetting); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	if err := ctx.BodyParser(&notificationSetting); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Omit("deleted_at", "created_at").Save(&notificationSetting); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(notificationSetting)

}
