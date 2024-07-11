package NET

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	db "worker/config"
	query "worker/config"
	"worker/controller/CLI"
	"worker/models"
)

var sqlGetLinserRoll string = query.ApgwXmlParsher("getLinserRoll")

func ApgwFrontendInfo(c *fiber.Ctx) error {
	db.Connection()
	output := CLI.SteampipeQuery(sqlGetLinserRoll)

	var linserRollResponse models.LinserRollResponse
	var frontendResource models.FrontendResource
	var linserRollResualt models.LinserRollResualt

	if err := json.Unmarshal([]byte(output), &linserRollResponse); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	for _, rowData := range linserRollResponse.Rows {
		if err := db.DB.Select("port").Where("frontend_id = ?", rowData.PortID).First(&frontendResource).Error; err != nil {
			log.Printf("Failed to find frontend resource: %s", err)
		} else {
			fmt.Println(rowData)
			fmt.Println(frontendResource)
			//continue
		}
		rowRemake := models.LinserRollReplace{
			ApgwName:   rowData.ApgwName,
			Hosts:      rowData.Hosts,
			LinserName: rowData.LinserName,
			Port:       frontendResource.Port,
		}

		linserRollResualt.Rows = append(linserRollResualt.Rows, rowRemake)
	}

	//fmt.Println(linserRollResponse)
	tmp, _ := json.Marshal(linserRollResualt)
	fmt.Println(string(tmp))

	return c.JSON(linserRollResualt)
}
