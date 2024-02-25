package api

import "github.com/gofiber/fiber/v2"

func Internal(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}
func Forbidden(message string) error {
	return fiber.NewError(fiber.StatusForbidden, message)
}
func Bad(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}
func Respond(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(&fiber.Map{
		"data": data,
	})
}
