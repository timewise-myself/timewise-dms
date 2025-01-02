package common

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterHandler(router fiber.Router, db *gorm.DB, registerFunc func(handler Handler)) {
	handler := Handler{
		Router: router,
		DB:     db,
	}
	registerFunc(handler)
}
