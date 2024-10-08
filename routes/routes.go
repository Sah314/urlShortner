package routes

import (
	"github.com/gofiber/fiber/v2"
	"urlShortner/helpers/middlewares"
	"urlShortner/services"
)

func SetupRoutes(app *fiber.App) {

	//app.Get("/:url", middlewares.ValidateJwt, services.ResolveURL)
	app.Get("/:url", services.ResolveURL)
	app.Post("api/v1", middlewares.ValidateJwt, services.ShortenURL)
	//app.Post("api/v1", services.ShortenURL)

	app.Post("api/v1/signup", services.Signup)
	app.Post("api/v1/login", services.Login)
	app.Patch("api/v1/update", middlewares.ValidateJwt, services.UpdateUser)
}
