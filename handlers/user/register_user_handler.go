package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterUserHandler(router fiber.Router, db *gorm.DB) {
	userHandler := UserHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", userHandler.getUsers)
	router.Get("/get", userHandler.getUserByEmail)
	router.Get("/:user_id", userHandler.getUserById)
	router.Post("/", userHandler.createUser)
	router.Put("/:user_id", userHandler.updateUser)
	router.Delete("/:user_id", userHandler.deleteUser)
	router.Post("/get-create", userHandler.getOrCreateUser)
}
