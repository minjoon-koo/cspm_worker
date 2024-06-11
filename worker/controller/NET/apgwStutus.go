package NET

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os/exec"
	db "worker/config"
	query "worker/config"
	"worker/models"
)

var sqlAzCliAppGatewayHealthCheck string = query.ApgwXmlParsher("azCliAppGatewayHealthCheck")
var sqlGetAppGateway string = query.ApgwXmlParsher("getAppGateway")

func AzCliHealthCheck(c *fiber.Ctx) error {
	var backendPool []models.BackendPool
	/*
		cmd := exec.Command("sh", "-c", sqlAzCliAppGatewayHealthCheck) // Ensure this command is appropriate for your OS and SQL
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Command execution failed: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to execute health check")
		}
		log.Printf("Command output: %s", output)
	*/

	/*if err := db.DB.Select("app_gateway_name").Find(&backendPool).Error; err != nil {
		log.Printf("Failed to query database: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Database query failed")
	}*/
	db.DB.Select("app_gateway_name").Find(&backendPool)

	/*
		backendPool 중복 항목 제거하여 리스트 새로 만들기
		ex)
		[0] ==> dev-apgw
		[1] ==> product-apgw
	*/
	appGatewayMap := make(map[string]bool)
	var uniqueBackendPool []models.BackendPool

	for _, pool := range backendPool {
		if _, exists := appGatewayMap[pool.AppGatewayName]; !exists {
			appGatewayMap[pool.AppGatewayName] = true
			uniqueBackendPool = append(uniqueBackendPool, pool)
		}
	}

	for _, pool := range uniqueBackendPool {
		cmd := exec.Command("sh", "-c", sqlAzCliAppGatewayHealthCheck+pool.AppGatewayName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Command execution failed: %v", err)
		}
		return c.SendString(string(output))
	}

	return c.Status(200).JSON(fiber.Map{
		"data": uniqueBackendPool[0],
	})

	//return c.SendString("a")
}
