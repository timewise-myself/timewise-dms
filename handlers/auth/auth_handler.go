package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"github.com/timewise-team/timewise-models/models"
	"time"
)

func (h *AuthHandler) CreateNewUser(c *fiber.Ctx) error {
	var registerResponseDto user_register_dto.RegisterResponseDto
	if err := c.BodyParser(&registerResponseDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user := models.TwUser{
		FirstName:   registerResponseDto.FirstName,
		LastName:    registerResponseDto.LastName,
		Email:       registerResponseDto.Email,
		LastLoginAt: time.Now(),
		Role:        "user",
	}
	if result := h.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON("register successfully")
}
