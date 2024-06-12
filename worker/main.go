package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	db "worker/config"
	"worker/routes"
)

func main() {
	godotenv.Load(".env")
	//export_port := os.Getenv("EXPORT_PORT")
	fmt.Println("Connecting to DB...")
	db.Connection()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(logger.New())
	routes.Setup(app)
	//controller.GetAllUserInfo()

	app.Listen("0.0.0.0:30001")

}
