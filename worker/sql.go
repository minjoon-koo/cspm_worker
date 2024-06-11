package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	db "worker/config"
	"worker/models"
)

var sql string = "select display_name, id , description, created_date_time, expiration_date_time from azuread_group"

var sql2 string = "select display_name, mail, department, member_of from azuread_user"

var azcmd string = "az network application-gateway show-backend-health --resource-group rg-soldout --name "

func main() {
	db.Connection()

	info()
	for _, status := range backendStatuses {
		fmt.Printf("AppGateway: %s\n", status.ApgwName)
		for _, statue := range status.Statues {
			fmt.Printf("  Address: %s, Health: %s\n", statue.Address, statue.Health)
		}
	}
}

var backendStatuses []models.BackendStatus

func info() {
	var backendPool []models.BackendPool
	db.DB.Select("app_gateway_name").Find(&backendPool)

	appGatewayMap := make(map[string]bool)
	var uniqueBackendPool []models.BackendPool

	//var backendStatuses []models.BackendStatus

	for _, pool := range backendPool {
		if _, exists := appGatewayMap[pool.AppGatewayName]; !exists {
			appGatewayMap[pool.AppGatewayName] = true
			uniqueBackendPool = append(uniqueBackendPool, pool)
		}
	}

	for _, pool := range uniqueBackendPool {
		var statues []models.StatuesResponse
		fmt.Println(pool.AppGatewayName)
		cmd := exec.Command("sh", "-c", azcmd+" "+pool.AppGatewayName)
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
					fmt.Printf("Address: %s, Health: %s\n", address, health)
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
	for _, status := range backendStatuses {
		fmt.Printf("AppGateway: %s\n", status.ApgwName)
		for _, statue := range status.Statues {
			fmt.Printf("  Address: %s, Health: %s\n", statue.Address, statue.Health)
		}
	}

}
