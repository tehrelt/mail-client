package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/lib"
)

func CheckPop3Auth(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		if app.User == nil {
			return Forbidden(ErrNotAuthenticated.Error())
		}

		connection, err := lib.Pop3Auth(app.Config.Pop3, app.User)
		if err != nil {
			return Internal(fmt.Sprintf("CheckPop3Auth [middleware]: %s", err.Error()))
		}

		ctx.Locals("connection", connection)

		return ctx.Next()
	}
}

func CheckSmtpAuth(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		if app.User == nil {
			return Forbidden(ErrNotAuthenticated.Error())
		}

		connection, err := lib.SmtpAuth(app.Config.Smtp, app.User)
		if err != nil {
			return Internal(fmt.Sprintf("CheckSmtpAuth [middleware]: %s", err.Error()))
		}

		ctx.Locals("connection", connection)

		return ctx.Next()
	}
}
