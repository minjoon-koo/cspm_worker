package NET

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	"log"
	"os/exec"
	db "worker/config"
	query "worker/config"
	"worker/models"
)

var sqlAppgatewayBeckend string = query.ApgwXmlParsher("getBackendName")

func UpdateAppgatewayBeckend() {
	var backendPoolResponse models.BackendPollResponse
	cmd := exec.Command("steampipe", "query", sqlAppgatewayBeckend, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute query: %s", string(output))
	}

	if err := json.Unmarshal(output, &backendPoolResponse); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	var existingBackendPool []models.BackendPool
	db.DB.Find(&existingBackendPool)
	existingMap := make(map[string]models.BackendPool)
	for _, backend := range existingBackendPool {
		existingMap[backend.BackendID] = backend
	}

	for _, backend := range backendPoolResponse.Rows {
		upsertBackendPool(backend)
		delete(existingMap, backend.BackendID)
	}

	for backendID := range existingMap {
		deleteBackendPool(backendID)
	}
}

func upsertBackendPool(backend models.BackendPool) {
	result := db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "backend_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"app_gateway_name": backend.AppGatewayName,
			"backend_name":     backend.BackendName,
			"ipaddress":        backend.Ipaddress,
			"backend_id":       backend.BackendID,
		}),
	}).Create(&backend)

	if result.Error != nil {
		log.Fatalf("Failed to upsert backend pool: %s", result.Error)
	} else {
		log.Println("Backend pool upserted successfully")
	}
}

func deleteBackendPool(backendID string) {
	result := db.DB.Where("backend_id = ?", backendID).Delete(&models.BackendPool{})
	if result.Error != nil {
		log.Fatalf("Failed to delete backend pool: %s", result.Error)
	} else {
		log.Println("Backend pool deleted successfully: %s", backendID)
	}
}
