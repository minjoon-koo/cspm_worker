package controller

import (
	"github.com/gofiber/fiber/v2"
	"worker/controller/HTTP"
	"worker/controller/IAM"
)

var sql string = "select display_name, id , description, created_date_time, expiration_date_time from azuread_group"

func GetAllUserInfo(c *fiber.Ctx) error { //그룹 업데이트
	IAM.UpdateAdGroups()
	IAM.UpdateDirectoryRole()
	steamSQLString := HTTP.SendGet("/IAM/user/worker")
	//HTTP.SendPost("/IAM/user/GF-api", steamSQLString)
	//println(steamSQLString)
	result, _ := HTTP.SendPost("/IAM/user/GF-api", steamSQLString)
	return c.SendString(string(result))
}
