package notification_setting

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterNotificationSettingHandler(router fiber.Router, db *gorm.DB) {
	notificationSettingHandler := NotificationSettingHandler{
		DB: db,
	}
	router.Get("/:user_id", notificationSettingHandler.GetNotificationSettingByUserId)
	router.Post("/", notificationSettingHandler.CreateNotificationSetting)
	router.Put("/:user_id", notificationSettingHandler.UpdateNotificationSetting)

}
