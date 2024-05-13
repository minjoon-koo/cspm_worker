package IAM

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	"log"
	"os/exec"
	db "worker/config"
	query "worker/config"
	"worker/models"
)

var sqlADGroup string = query.XmlParsher("getAllGroup")
var sqlDirectoryRole string = query.XmlParsher("getDirectory_Role")

/*
AD Group 파싱
*/
func UpdateAdGroups() {
	var adGroupResponse models.AdGroupResponse
	// Steampipe 쿼리 실행
	cmd := exec.Command("steampipe", "query", sqlADGroup, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Steampipe query failed: %s", string(output))
	}

	// JSON 파싱
	//var response models.AdGroupResponse
	if err := json.Unmarshal(output, &adGroupResponse); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	// 현재 DB의 상태를 캐시
	var existingGroups []models.AdGroup
	db.DB.Find(&existingGroups)
	existingMap := make(map[string]models.AdGroup)
	for _, group := range existingGroups {
		existingMap[group.Id] = group
	}

	// 데이터베이스에 업데이트/삽입
	for _, group := range adGroupResponse.Rows {
		upsertAdGroup(group)
		delete(existingMap, group.Id)

	}
	for id := range existingMap {
		deleteGroup(id)
	}
}

/*
Directory Role
*/
func UpdateDirectoryRole() {
	var directoryRoleResponse models.DirectoryRoleResponse
	cmd := exec.Command("steampipe", "query", sqlDirectoryRole, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Steampipe query failed: %s", string(output))
	}

	if err := json.Unmarshal(output, &directoryRoleResponse); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	var existingGroups []models.DirectoryRole
	db.DB.Find(&existingGroups)
	existingMap := make(map[string]models.DirectoryRole)
	for _, group := range existingGroups {
		existingMap[group.Id] = group
	}
	for _, role := range directoryRoleResponse.Rows {
		//upsertAdGroup(group)
		upsertDirectoryRole(role)
		delete(existingMap, role.Id)
	}
	for id := range existingMap {
		deleteDirectoryRole(id)
	}
}

/*
SQL
*/
//AD Group SQL func
func upsertAdGroup(group models.AdGroup) {
	// Upsert 로직: ID가 동일할 경우 다른 필드 업데이트
	result := db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // ID를 고유 키로 사용
		DoUpdates: clause.Assignments(map[string]interface{}{
			"display_name":         group.DisplayName,
			"description":          group.Description,
			"create_date_time":     group.CreateDateTime,
			"expiration_date_time": group.ExpirationDateTime,
		}),
	}).Create(&group)

	if result.Error != nil {
		log.Printf("Failed to upsert ad group: %s", result.Error)
	} else {
		log.Println("Ad group upserted successfully")
	}
}

func deleteGroup(id string) {
	// 삭제 로직
	result := db.DB.Where("id = ?", id).Delete(&models.AdGroup{})
	if result.Error != nil {
		log.Printf("Failed to delete ad group: %s", result.Error)
	} else {
		log.Printf("Ad group deleted successfully: %s", id)
	}
}

// Directory Role SQL func
func upsertDirectoryRole(role models.DirectoryRole) {
	result := db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"display_name": role.Display_name,
		}),
	}).Create(&role)
	if result.Error != nil {
		log.Printf("Failed to upsert directory Role: %s", result.Error)
	} else {
		log.Println("Directory Role upserted successfully")
	}
}

func deleteDirectoryRole(id string) {
	result := db.DB.Where("id = ?", id).Delete(&models.DirectoryRole{})
	if result.Error != nil {
		log.Printf("Failed to delete directory Role: %s", result.Error)
	} else {
		log.Printf("directory Role deleted successfully: %s", id)
	}
}
