package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-api/internal/client"
	"go-api/internal/handler"
	"go-api/internal/middleware"
)

func main() {
	nodeURL := os.Getenv("NODE_API_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	nodeClient := client.NewNodeClient(nodeURL, jwtSecret)
	matrixHandler := handler.NewMatrixHandler(nodeClient)

	app := fiber.New(fiber.Config{
		BodyLimit: 4 * 1024 * 1024,
	})

	app.Use(logger.New())
	middleware.SetupSecurityMiddlewares(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "go-api"})
	})

	api := app.Group("/api", middleware.KeyAuth())
	api.Post("/factorization", matrixHandler.ProcessMatrix)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
