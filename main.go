package main

import (
	"log"

	cassandra "github.com/HarishGangula/content-state-go-concurrency/cassandra"
	models "github.com/HarishGangula/content-state-go-concurrency/models"
	"github.com/gofiber/fiber/v2"
)

func startServer() {
	app := fiber.New()

	app.Patch("/api/v1/contentstate/update", func(c *fiber.Ctx) error {

		request := new(models.Request)

		if err := c.BodyParser(request); err != nil {
			log.Fatal(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})

		} else {
			responseChan := make(chan models.Response)
			go cassandra.UpsertContentState(*request, responseChan)
			response := <-responseChan
			return c.Status(fiber.StatusOK).JSON(response)
		}
	})

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}

func main() {

	cassandra.Init()
	startServer()
	defer cassandra.Close()
}
