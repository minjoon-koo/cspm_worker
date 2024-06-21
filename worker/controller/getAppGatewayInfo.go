package controller

import (
	"github.com/gofiber/fiber/v2"
	"worker/controller/NET"
)

func GetAppGatewayInfo(c *fiber.Ctx) error { //그룹 업데이트

	NET.UpdateAppgatewayBeckend()
	//return c.SendString("{\"url\":\"https://soldout.co.kr\",\"num\":\"80\"}")
	return c.SendString("{\"num\": {\"WAF1\":\"80\"}}")
	//return c.SendString("1,2,3,4,5,6")
}
