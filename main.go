package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: time.Second * 5,
		WriteTimeout: time.Second * 5,
		ReadTimeout: time.Second * 5,
	})

	// Routing
	app.Get("/", func (ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World!")
	})

	// Start application
	err := app.Listen("localhost:3000")

	if (err != nil) {
		panic(err)
	}
}