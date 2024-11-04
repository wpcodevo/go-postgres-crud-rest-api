package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/go-postgres-crud-rest-api/initializers"
)

func main() {
	env, err := initializers.LoadEnv(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&env)
	app := fiber.New()

	app.Get("/api/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "CRUD Operations on PostgreSQL using Golang REST API",
		})
	})

	log.Fatal(app.Listen(":" + env.ServerPort))
}
