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
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new Fiber instance
	app := fiber.New()

	// Use logger middleware
	app.Use(logger.New())

	// Setup routes (your custom route setup logic)
	setupRoutes(app)

	// Get the port from environment variables
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":3000" // Default port if not set
	}

	// Log the message before starting the server
	log.Printf("Server running on port %s", port)

	// Start the Fiber server on the specified port
	log.Fatal(app.Listen(port))
}
