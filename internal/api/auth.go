package api

import (
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/dto"
	"mail-client/internal/lib"
)

func HandlerAuth(api *API) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		if api.User != nil {
			return ctx.SendStatus(204)
		}

		var req dto.User

		if err := ctx.BodyParser(&req); err != nil {
			return Forbidden(err.Error())
		}

		connection, err := lib.Pop3Auth(api.Config.Pop3, &req)
		if err != nil {
			return Internal(err.Error())
		}

		defer connection.Quit()

		api.User = &req

		return ctx.SendStatus(200)
	}
}

func HandlerAuthAlive(api *API) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := api.User

		if user == nil {
			return Bad("not authenticated")
		}

		connection, err := lib.Pop3Auth(api.Config.Pop3, user)
		if err != nil {
			return Forbidden(err.Error())
		}

		defer connection.Quit()

		return Respond(ctx, fiber.Map{
			"user": user.User,
		})
	}
}

func HandlerLogout(api *API) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		api.User = nil

		return ctx.SendStatus(200)
	}
}
