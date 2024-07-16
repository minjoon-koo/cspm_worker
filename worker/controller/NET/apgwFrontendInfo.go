package NET

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	db "worker/config"
	query "worker/config"
	"worker/controller/CLI"
	"worker/controller/HTTP"
	"worker/models"
)

var sqlGetLinserRoll string = query.ApgwXmlParsher("getLinserRoll")

func ApgwFrontendInfo(c *fiber.Ctx) error {
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

	tmp, _ := json.Marshal(linserRollResualt)
	fmt.Println(string(tmp))

	return c.JSON(linserRollResualt)
}

func ApgwFrontendPort(c *fiber.Ctx) error {
	totalResault := HTTP.SendGet("/NET/GW/Info")
	linserId := c.Params("linserId")

	/*
		totalResult 값
		{"rows":[{"apgw_name":"production-apgw","hosts":"[\"standby.soldout.co.kr\", \"standby*.soldout.co.kr\", \"*standby.soldout.co.kr\"]","linser_name":"Listner-Standby-HTTPS","port":"443"},{"apgw_name":"production-apgw","hosts":"[\"standby.soldout.co.kr\", \"standby*.soldout.co.kr\", \"*standby.soldout.co.kr\"]","linser_name":"Listner-Standby-HTTP","port":"80"},{"apgw_name":"production-apgw","hosts":"[\"soldout.co.kr\", \"*.soldout.co.kr\"]","linser_name":"Listner-HTTP-Production01","port":"80"},{"apgw_name":"production-apgw","hosts":"[\"soldout.co.kr\", \"www.soldout.co.kr\", \"*.soldout.co.kr\"]","linser_name":"Listner-HTTPS-Productiono01","port":"443"},{"apgw_name":"dev-apgw","hosts":"[\"test-*.soldout.co.kr\", \"api-test.soldout.co.kr\", \"dev-*.soldout.co.kr\", \"sk-lb.soldout.co.kr\", \"test.soldout.co.kr\"]","linser_name":"listner-test01-01","port":"443"},{"apgw_name":"dev-apgw","hosts":"[\"t2.soldout.co.kr\", \"t2-*.soldout.co.kr\", \"api-t2.soldout.co.kr\"]","linser_name":"listner-t2-01","port":"443"},{"apgw_name":"dev-apgw","hosts":"[\"qa.soldout.co.kr\", \"qa-*.soldout.co.kr\", \"*qa.soldout.co.kr\", \"api-dev.soldout.co.kr\"]","linser_name":"listner-qa-01","port":"443"},{"apgw_name":"dev-apgw","hosts":"[\"t2.soldout.co.kr\", \"t2-*.soldout.co.kr\", \"api-t2.soldout.co.kr\"]","linser_name":"listner-t2-redirect","port":"80"},{"apgw_name":"dev-apgw","hosts":"[\"qa.soldout.co.kr\", \"qa-*.soldout.co.kr\", \"*qa.soldout.co.kr\", \"api-dev.soldout.co.kr\"]","linser_name":"listner-qa-redirect","port":"80"},{"apgw_name":"dev-apgw","hosts":"[\"stage.soldout.co.kr\", \"*-stage.soldout.co.kr\", \"stage-*.soldout.co.kr\"]","linser_name":"listner-stage-01","port":"443"},{"apgw_name":"dev-apgw","hosts":"[\"stage.soldout.co.kr\", \"*-stage.soldout.co.kr\", \"stage-*.soldout.co.kr\"]","linser_name":"listner-stage-redirect","port":"80"},{"apgw_name":"dev-apgw","hosts":"[\"api-t1.soldout.co.kr\", \"t1.soldout.co.kr\", \"t1-*.soldout.co.kr\"]","linser_name":"listner-test09-redirect","port":"80"},{"apgw_name":"dev-apgw","hosts":"[\"test-*.soldout.co.kr\", \"api-test.soldout.co.kr\", \"dev-*.soldout.co.kr\", \"test.soldout.co.kr\"]","linser_name":"listner-test01-redirect","port":"80"},{"apgw_name":"dev-apgw","hosts":"[\"api-t1.soldout.co.kr\", \"t1.soldout.co.kr\", \"t1-*.soldout.co.kr\"]","linser_name":"listner-test09-01","port":"443"}]}
		linserId 값을 json에서 검색하고 port 번호를 뽑으려고 함
		예시) linserId = Listner-Standby-HTTPS  --> 결과 {"port":"443"}
	*/

	var selectResault models.Frontends
	if err := json.Unmarshal([]byte(totalResault), &selectResault); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
		return c.Status(500).JSON(fiber.Map{
			"port": "null",
		})
	}

	for _, row := range selectResault.Rows {
		if row.LinserName == linserId {
			return c.JSON(fiber.Map{
				"port": row.Port,
			})
		}
	}

	return c.JSON(fiber.Map{
		"port": "not found",
	})

}
