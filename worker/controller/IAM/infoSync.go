package IAM

import (
	"encoding/json"
	"log"
	"os/exec"
	"sync"
	db "worker/config"
	query "worker/config"
	"worker/models"

	"gorm.io/gorm/clause"
)

var sqlADGroup string = query.IamXmlParsher("getAllGroup")
var sqlDirectoryRole string = query.IamXmlParsher("getDirectory_Role")

func runSteampipeQuery(query string, wg *sync.WaitGroup, resultChan chan<- []byte, errChan chan<- error) {
	defer wg.Done()
	cmd := exec.Command("steampipe", "query", query, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		errChan <- err
		return
	}
	resultChan <- output
}

/*
AD Group 파싱
*/
func UpdateAdGroups() {
	var wg sync.WaitGroup
	resultChan := make(chan []byte)
	errChan := make(chan error)

	wg.Add(1)
	go runSteampipeQuery(sqlADGroup, &wg, resultChan, errChan)

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	select {
	case output := <-resultChan:
		var adGroupResponse models.AdGroupResponse
		if err := json.Unmarshal(output, &adGroupResponse); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %s", err)
		}

		var existingGroups []models.AdGroup
		db.DB.Find(&existingGroups)
		existingMap := make(map[string]models.AdGroup)
		for _, group := range existingGroups {
			existingMap[group.Id] = group
		}

		for _, group := range adGroupResponse.Rows {
			upsertAdGroup(group)
			delete(existingMap, group.Id)
		}
		for id := range existingMap {
			deleteGroup(id)
		}
	case err := <-errChan:
		log.Fatalf("Steampipe query failed: %s", err)
	}
}

/*
Directory Role
*/
func UpdateDirectoryRole() {
	var wg sync.WaitGroup
	resultChan := make(chan []byte)
	errChan := make(chan error)

	wg.Add(1)
	go runSteampipeQuery(sqlDirectoryRole, &wg, resultChan, errChan)

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	select {
	case output := <-resultChan:
		var directoryRoleResponse models.DirectoryRoleResponse
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
			upsertDirectoryRole(role)
			delete(existingMap, role.Id)
		}
		for id := range existingMap {
			deleteDirectoryRole(id)
		}
	case err := <-errChan:
		log.Fatalf("Steampipe query failed: %s", err)
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
