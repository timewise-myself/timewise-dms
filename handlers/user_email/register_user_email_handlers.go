package user_email

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserEmailHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterUserEmailHandler(router fiber.Router, db *gorm.DB) {
	userEmailHandler := UserEmailHandler{
		Router: router,
		DB:     db,
	}

	router.Get("/", userEmailHandler.getUserEmails)
	router.Get("/user/:user_id", userEmailHandler.getUserEmailByUserId)
	router.Get("/email/:email", userEmailHandler.getUserEmailByEmail)
	router.Get("/check", userEmailHandler.getUserEmailToCheckBeforeLink)
	router.Post("/", userEmailHandler.createUserEmail)
	router.Patch("/", userEmailHandler.updateStatusInUserEmail)
	router.Patch("/status", userEmailHandler.updateUserEmailStatusAndIsLinkedTo)
	router.Delete("/", userEmailHandler.deleteUserEmail)
	router.Get("/search/:query", userEmailHandler.searchUserEmail)
	router.Get("/listApprove/:scheduleId", userEmailHandler.getEmailInProgress)
	router.Get("/clear-expired", userEmailHandler.clearExpiredUserEmails)
	router.Get("/user_id/:user_id", userEmailHandler.getExactUserEmailByUserId)
	router.Get("clear-rejected", userEmailHandler.clearStatusRejectedEmail)
}
