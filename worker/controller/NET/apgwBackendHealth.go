package NET

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"os/exec"
	db "worker/config"
	query "worker/config"
	"worker/models"
)

// zacmd string
var sqlAzCliAppGatewayHealthCheck string = query.ApgwXmlParsher("azCliAppGatewayHealthCheck")
var sqlGetAppGateway string = query.ApgwXmlParsher("getAppGateway")
var azcmd string = "az network application-gateway show-backend-health --resource-group rg-soldout --name "

/*
health check
1. db query > sync 데이터 확인
2. cli 요청 > health 상태 확인 > 매핑
3. 상태 값 리턴
*/
var backendStatuses []models.BackendStatus

func HealthInfo(c *fiber.Ctx) error {
	var backendPool []models.BackendPool
	db.DB.Select("app_gateway_name").Find(&backendPool)
	appGatewayMap := make(map[string]bool)

	var uniqueBackendPool []models.BackendPool
	for _, pool := range backendPool {
		if _, exists := appGatewayMap[pool.AppGatewayName]; !exists {
			appGatewayMap[pool.AppGatewayName] = true
			uniqueBackendPool = append(uniqueBackendPool, pool)
		}
	}

	for _, pool := range uniqueBackendPool {
		var statues []models.StatuesResponse
		//fmt.Println(pool.AppGatewayName)
		cmd := exec.Command("sh", "-c", sqlAzCliAppGatewayHealthCheck+" "+pool.AppGatewayName)
		output, err := cmd.CombinedOutput()
		jsonData := string(output)

		//fmt.Println("Command Output:", jsonData)
		var result interface{}
		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			log.Printf("JSON Unmarshal failed: %v", err)
			continue
		}
		// 동적으로 파싱된 JSON 데이터 출력
		//fmt.Printf("Parsed JSON: %+v\n", result)

		jsonMap, ok := result.(map[string]interface{})
		if !ok {
			log.Printf("Failed to convert result to map[string]interface{}")
			continue
		}

		backendAddressPools, ok := jsonMap["backendAddressPools"].([]interface{})
		if !ok {
			log.Printf("Failed to get backendAddressPools")
			continue
		}

		for _, pool2 := range backendAddressPools {
			poolMap, ok := pool2.(map[string]interface{})
			if !ok {
				continue
			}
			backendHttpSettingsCollection, ok := poolMap["backendHttpSettingsCollection"].([]interface{})
			if !ok {
				continue
			}

			for _, settings := range backendHttpSettingsCollection {
				settingsMap, ok := settings.(map[string]interface{})
				if !ok {
					continue
				}
				servers, ok := settingsMap["servers"].([]interface{})
				if !ok {
					continue
				}
				for _, server := range servers {
					serverMap, ok := server.(map[string]interface{})
					if !ok {
						continue
					}
					address, _ := serverMap["address"].(string)
					health, _ := serverMap["health"].(string)
					//fmt.Printf("Address: %s, Health: %s\n", address, health)
					statue := models.StatuesResponse{
						Address: address,
						Health:  health,
					}
					statues = append(statues, statue)
				}
			}
		}
		backendStatus := models.BackendStatus{
			ApgwName: pool.AppGatewayName,
			Statues:  statues,
		}
		backendStatuses = append(backendStatuses, backendStatus)
	}

	jsonOutput, err := json.Marshal((backendStatuses))
	if err != nil {
		log.Printf("JSON Marshal failed: %v", err)
	}
	jsonString := string(jsonOutput)

	return c.SendString(jsonString)
}
