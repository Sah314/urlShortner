package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
	"urlShortner/services"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", services.ResolveURL)
	app.Post("api/v1", services.ShortenURL)

}

func main() {
	err := godotenv.Load()
	if err != nil {

	}
	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}
