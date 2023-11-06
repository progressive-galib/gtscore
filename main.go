package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	SetupRoutes(app)
	port := 8080
	fmt.Printf("Listening on :%d\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
