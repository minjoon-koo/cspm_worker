package routes

import (
	"github.com/gofiber/fiber/v2"
	"worker/controller"
	"worker/controller/IAM"
)

func Setup(app *fiber.App) {
	app.Get("/IAM/user/worker", IAM.SteamQLUserGet)
	app.Post("/IAM/user/GF-api", IAM.ResultUserGet)
	app.Get("/IAM/", controller.GetAllUserInfo)
}
