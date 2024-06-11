package routes

import (
	"github.com/gofiber/fiber/v2"
	"worker/controller"
	"worker/controller/IAM"
	"worker/controller/NET"
)

func Setup(app *fiber.App) {
	app.Get("/IAM/user/worker", IAM.SteamQLUserGet)
	app.Post("/IAM/user/GF-api", IAM.ResultUserGet)
	app.Get("/IAM/", controller.GetAllUserInfo)

	app.Get("/NET/GW", controller.GetAppGatewayInfo)
	app.Get("/NET/apgw/health", NET.AzCliHealthCheck)
}
