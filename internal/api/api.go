package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"mail-client/internal/config"
	"mail-client/internal/dto"
)

type API struct {
	app    *fiber.App
	User   *dto.User
	Config *config.AppConfig
}

func Start(cfg *config.AppConfig) error {
	app := fiber.New(fiber.Config{
		AppName:       "mail-client-api",
		CaseSensitive: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = ctx.Status(code).JSON(fiber.Map{
				"message": e.Message,
			})

			return nil
		},
	})

	api := &API{
		app:    app,
		Config: cfg,
	}

	api.app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))
	api.app.Use(logger.New())

	api.configure()

	return api.app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (api *API) configure() {

	api.app.Get("/auth", HandlerAuthAlive(api))
	api.app.Post("/auth", HandlerAuth(api))
	api.app.Post("/logout", HandlerLogout(api))

	pop3 := api.app.Group("/pop3", CheckPop3Auth(api))
	pop3.Get("/list", HandlerPopList(api))
	pop3.Get("/list/:id", HandlerPopRetrieve(api))
	pop3.Get("/stat", HandlerPopStat(api))

	smtp := api.app.Group("/smtp", CheckSmtpAuth(api))
	smtp.Post("/send", HandlerSmtpSend(api))
}
