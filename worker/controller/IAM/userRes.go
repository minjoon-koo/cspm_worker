package IAM

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os/exec"
	db "worker/config"
	query "worker/config"
	"worker/models"
)

var sql string = query.IamXmlParsher("getAllUsers")

func SteamQLUserGet(c *fiber.Ctx) error {
	cmd := exec.Command("steampipe", "query", sql, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Steampipe query failed: %s", string(output))
	}

	//HTTP.SendPost("/IAM/user/GF-api", string(output))

	return c.SendString(string(output))
}

func ResultUserGet(c *fiber.Ctx) error {
	var AdGroup models.AdGroup
	var DirectoryRole models.DirectoryRole
	var steamQLResult models.AzureADUserResponse
	if err := c.BodyParser(&steamQLResult); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	replaceResult := models.AzureADuserRequest{
		Rows: make([]models.AzureADUserRemake, 0, len(steamQLResult.Rows)),
	}

	for _, rowData := range steamQLResult.Rows {
		memberDisplayNames := make([]string, 0) // 문자열 슬라이스로 변경
		for _, gid := range rowData.MemberOf {
			displayName := "" // 초기 문자열 선언
			if gid.OdataType == "#microsoft.graph.group" {
				if err := db.DB.Select("display_name").Where("id = ?", gid.ID).First(&AdGroup).Error; err != nil {
					log.Printf("Error retrieving group info: %v", err)
					continue
				}
				displayName = AdGroup.DisplayName
			} else if gid.OdataType == "#microsoft.graph.directoryRole" {
				if err := db.DB.Select("display_name").Where("id = ?", gid.ID).First(&DirectoryRole).Error; err != nil {
					log.Printf("Error retrieving directory role info: %v", err)
					continue
				}
				displayName = DirectoryRole.Display_name
			}
			if displayName != "" {
				memberDisplayNames = append(memberDisplayNames, displayName)
			}
		}

		userRemake := models.AzureADUserRemake{
			DisplayName: rowData.DisplayName,
			Mail:        rowData.Mail,
			Department:  rowData.Department,
			MemberOf:    memberDisplayNames, // 문자열 슬라이스 할당
		}
		replaceResult.Rows = append(replaceResult.Rows, userRemake)
	}

	// JSON 결과를 클라이언트에 반환
	return c.JSON(replaceResult)
}

/*
func ResultUserGet2(c *fiber.Ctx) error {
	var AdGroup models.AdGroup
	var DirectoryRole models.DirectoryRole
	var steamQLResult models.AzureADUserResponse
	if err := c.BodyParser(&steamQLResult); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	replaceResult := models.AzureADuserRequest{
		Rows: make([]models.AzureADUserRemake, 0, len(steamQLResult.Rows)),
	}

	for _, rowData := range steamQLResult.Rows {
		memberNames := make([]models.GroupName, 0)
		for _, gid := range rowData.MemberOf {
			groupName := models.GroupName{gid.ID}
			if gid.OdataType == "#microsoft.graph.group" {
				if err := db.DB.Select("display_name").Where("id = ?", gid.ID).First(&AdGroup).Error; err != nil {
					log.Printf("Error retrieving group info: %v", err)
					continue
				}
				groupName.GroupName = AdGroup.DisplayName
			} else if gid.OdataType == "#microsoft.graph.directoryRole" {
				if err := db.DB.Select("display_name").Where("id = ?", gid.ID).First(&DirectoryRole).Error; err != nil {
					log.Printf("Error retrieving directory role info: %v", err)
					continue
				}
				groupName.GroupName = DirectoryRole.Display_name
			}
			memberNames = append(memberNames, groupName)
		}

		userRemake := models.AzureADUserRemake{
			DisplayName: rowData.DisplayName,
			Mail:        rowData.Mail,
			Department:  rowData.Department,
			MemberOf:    memberNames,
		}
		replaceResult.Rows = append(replaceResult.Rows, userRemake)
	}

	// JSON 결과를 클라이언트에 반환
	return c.JSON(replaceResult)
}*/
