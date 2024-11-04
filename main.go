package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wpcodevo/go-postgres-crud-rest-api/controllers"
	"github.com/wpcodevo/go-postgres-crud-rest-api/initializers"
)

func main() {
	env, err := initializers.LoadEnv(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&env)
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	micro.Route("/feedbacks", func(router fiber.Router) {
		router.Post("/", controllers.CreateFeedbackHandler)
		router.Get("", controllers.FindFeedbacksHandler)
	})

	micro.Route("/feedbacks/:feedbackId", func(router fiber.Router) {
		router.Get("", controllers.FindFeedbackByIdHandler)
		router.Patch("", controllers.UpdateFeedbackHandler)
		router.Delete("", controllers.DeleteFeedbackHandler)
	})

	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "CRUD Operations on PostgreSQL using Golang REST API",
		})
	})

	log.Fatal(app.Listen(":" + env.ServerPort))
}
