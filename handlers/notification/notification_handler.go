package notification

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	DB *gorm.DB
}

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new notification
// @Tags notification
// @Accept json
// @Produce json
// @Param createNotificationRequest body core_dtos.PushNotificationDto true "Create notification request"
// @Success 200 {object} models.TwNotifications
// @Router /dbms/v1/notification [post]
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	// Get data from request
	var request models.TwNotifications
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if result := h.DB.Create(&request); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	// Insert data into database
	return nil
}

// GetUnsentNotifications godoc
// @Summary Get unsent notifications
// @Description Get unsent notifications
// @Tags notification
// @Accept json
// @Produce json
// @Success 200 {array} models.TwNotifications
// @Router /dbms/v1/notification [get]
func (h *NotificationHandler) GetUnsentNotifications(ctx *fiber.Ctx) error {
	var notifications []models.TwNotifications
	if err := h.DB.Where("is_sent = ?", false).Preload("UserEmail").Find(&notifications).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(notifications)
}

// GetNotiByUserEmailIds godoc
// @Summary Get notifications by user email ids
// @Description Get notifications by user email ids
// @Tags notification
// @Accept json
// @Produce json
// @Param user_email_ids body []string true "User email ids"
// @Success 200 {array} models.TwNotifications
// @Router /dbms/v1/notification/user-email-ids [post]
func (h *NotificationHandler) GetNotiByUserEmailIds(ctx *fiber.Ctx) error {
	var userEmailIds []string
	_ = ctx.BodyParser(&userEmailIds)
	var notifications []models.TwNotifications
	if err := h.DB.Where("user_email_id in (?)", userEmailIds).Preload("UserEmail").Find(&notifications).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(notifications)
}

// UpdateNotiStatus godoc
// @Summary Update notification status
// @Description Update notification status
// @Tags notification
// @Accept json
// @Produce json
// @Param notification_id query string true "Notification ID"
// @Param is_read query string true "Is read"
// @Success 200 {object} models.TwNotifications
// @Router /dbms/v1/notification/update-status/read [put]
func (h *NotificationHandler) UpdateNotiStatus(c *fiber.Ctx) error {
	notiId := c.Query("notification_id")
	if notiId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}
	isRead := c.Query("is_read")
	if isRead == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Is read is required",
		})
	}
	var notification models.TwNotifications
	if err := h.DB.First(&notification, notiId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notification not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if isRead == "true" {
		notification.IsRead = true
	} else {
		notification.IsRead = false
	}
	if err := h.DB.Save(&notification).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(notification)
}

// updateNotificationToSent godoc
// @Summary Update notification to sent
// @Description Update notification to sent
// @Tags notification
// @Accept json
// @Produce json
// @Param notification_id path string true "Notification ID"
// @Success 200 {object} fiber.Map
// @Router /dbms/v1/notification/{notification_id} [put]
func (h *NotificationHandler) updateNotificationToSent(ctx *fiber.Ctx) error {
	notificationID := ctx.Params("notification_id")
	if notificationID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	var notification models.TwNotifications
	if err := h.DB.First(&notification, notificationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notification not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	notification.IsSent = true
	if err := h.DB.Save(&notification).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Notification updated to sent",
	})
}
