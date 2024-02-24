package api

import "github.com/gofiber/fiber/v2"

func internal(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}
func forbidden(message string) error {
	return fiber.NewError(fiber.StatusForbidden, message)
}
func bad(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}
func respond(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(&fiber.Map{
		"data": data,
	})
}
