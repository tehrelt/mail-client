package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/lib"
)

func HandlerPopAuth(app *API) HandlerFunc {

	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	return func(ctx *fiber.Ctx) error {

		var req request

		if err := ctx.BodyParser(&req); err != nil {
			return internal("Internal server error: cannot parse json")
		}

		if err := app.pop.Auth(req.User, req.Password); err != nil {
			//if 0 == strings.Compare(err.Error(), "failed at USER command: something went wrong: -ERR Unknown command: USER") {
			//	return bad("Already authenticated")
			//}
			//
			//if strings.Contains(err.Error(), "-ERR Invalid login or password") {
			//	return forbidden("Invalid credentials")
			//}

			return internal(fmt.Sprintf("Cannot auth: %s", err))
		}

		return respond(ctx, "successfully auth")
	}
}

func HandlerPopList(app *API) HandlerFunc {

	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	type response struct {
		Messages []*lib.Mail `json:"messages"`
	}

	return func(ctx *fiber.Ctx) error {

		var req request
		var res response

		if err := ctx.BodyParser(&req); err != nil {
			return internal("Internal server error: cannot parse json")
		}

		listedMessages, err := app.pop.ListAll()
		if err != nil {
			return internal(fmt.Sprintf("pop.ListAll: %s", err))
		}

		for _, listedMessage := range listedMessages {
			msg, err := app.pop.Retrieve(listedMessage)
			if err != nil {
				return internal(fmt.Sprintf("pop.Retr: %s", err))
			}

			res.Messages = append(res.Messages, msg)
		}

		return respond(ctx, res)
	}
}
