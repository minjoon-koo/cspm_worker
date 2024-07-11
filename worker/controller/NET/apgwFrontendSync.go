package NET

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	"log"
	db "worker/config"
	query "worker/config"
	"worker/controller/CLI"
	"worker/models"
)

var sqlGetFrontendPort string = query.ApgwXmlParsher("getFrontendPort")

func UpdateFrontendResource() {
	var frontendResponse models.FrontendResponse
	output := CLI.SteampipeQuery(sqlGetFrontendPort)

	if err := json.Unmarshal(output, &frontendResponse); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	var existingfrontendResource []models.FrontendResource
	db.DB.Find(&existingfrontendResource)
	existingMap := make(map[string]models.FrontendResource)
	for _, frontend := range existingfrontendResource {
		existingMap[frontend.FrontendID] = frontend
	}

	for _, frontend := range frontendResponse.Rows {
		upsertFrontendResource(frontend)
		delete(existingMap, frontend.FrontendID)
	}

	for frontendID := range existingMap {
		deleteFrontendResource(frontendID)
	}
}

func upsertFrontendResource(frontend models.FrontendResource) {
	result := db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "frontend_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"apgw_name":     frontend.ApgwName,
			"frontend_name": frontend.FrontendName,
			"port":          frontend.Port,
			"frontend_id":   frontend.FrontendID,
		}),
	}).Create(&frontend)

	if result.Error != nil {
		log.Fatalf("Failed to upsert frontend resource: %s", result.Error)
	} else {
		log.Println("frontend resource upserted successfully")
	}
}

func deleteFrontendResource(frontendID string) {
	result := db.DB.Where("frontend_id = ?", frontendID).Delete(&models.FrontendResource{})
	if result.Error != nil {
		log.Fatalf("Failed to delete frontend resource: %s", result.Error)
	} else {
		log.Println("frontend resource deleted successfully: %s", frontendID)
	}
}
