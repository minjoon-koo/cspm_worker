package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	db "worker/config"
	"worker/routes"
)

func main() {
	fmt.Println("Connecting to DB...")
	db.Connection()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	routes.Setup(app)
	//controller.GetAllUserInfo()

	app.Listen(":3000")

}
