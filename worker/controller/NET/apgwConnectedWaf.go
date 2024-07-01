package NET

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	query "worker/config"
	"worker/controller/CLI"
)

var sqlGetFirewallPolicy string = query.ApgwXmlParsher("getFirewallPolicy")
var resourceUrl string = "https://portal.azure.com/#@x2soft.onmicrosoft.com/resource"

func ApgwConnectedWaf(c *fiber.Ctx) error {

	queryResult := CLI.SteampipeQuery(sqlGetFirewallPolicy)

	var result map[string]interface{}
	json.Unmarshal(queryResult, &result)

	rows, ok := result["rows"].([]interface{})
	if !ok {
		fmt.Println("No rows found")
	}

	dataMap := make(map[string]string)
	for _, row := range rows {
		rowMap, ok := row.(map[string]interface{})
		if !ok {
			fmt.Println("Error: row is not a map")
			continue
		}

		apgwName, ok := rowMap["apgw_name"].(string)
		if !ok {
			fmt.Println("Error: 'apgw_name' is not a string")
			continue
		}

		firewallPolicy, ok := rowMap["firewall_policy"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: 'firewall_policy' is not a map")
			continue
		}

		id, ok := firewallPolicy["id"].(string)
		if !ok {
			fmt.Println("Error: 'id' is not a string")
			continue
		}
		dataMap[apgwName] = id
	}

	newRows := make([]map[string]string, 0)
	for key, value := range dataMap {
		newRow := make(map[string]string)
		newRow["apgw_name"] = key
		if value == "null" {
			newRow["firewallPolicy"] = "Null"
			newRow["firewallPolicyName"] = "Null"
			newRow["resourceUri"] = "Null"
		} else {
			newRow["firewallPolicy"] = value
			parts := strings.Split(value, "/")
			newRow["firewallPolicyName"] = parts[len(parts)-1]
			newRow["resourceUri"] = resourceUrl + value
		}
		newRows = append(newRows, newRow)
	}

	finalResult := map[string]interface{}{
		"rows": newRows,
	}
	finalOutput, err := json.MarshalIndent(finalResult, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling final JSON:", err)
	}

	return c.SendString(string(finalOutput))
}
