package utils

import (
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/pkg/messages"
	"github.com/gofiber/fiber/v2"
)

func sendMessage(c *fiber.Ctx, status int, msg messages.Message) error {
	return c.Status(status).JSON(&model.ErrorResponse{Error: msg})
}

// 4**
func Send400(c *fiber.Ctx, msg messages.Message) error {
	return sendMessage(c, fiber.StatusBadRequest, msg)
}

func Send401(c *fiber.Ctx, msg messages.Message) error {
	return sendMessage(c, fiber.StatusUnauthorized, msg)
}

func Send404(c *fiber.Ctx, msg messages.Message) error {
	return sendMessage(c, fiber.StatusNotFound, msg)
}

// 5**
func Send500(c *fiber.Ctx, msg messages.Message) error {
	return sendMessage(c, fiber.StatusInternalServerError, msg)
}
