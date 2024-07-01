package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	query "worker/config"
)

var sql string = query.ApgwXmlParsher("getFirewallPolicy")
var url string = "https://portal.azure.com/#@x2soft.onmicrosoft.com/resource"

func main() {
	//db.Connection()
	cmd := exec.Command("steampipe", "query", sql, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(string(output))
	}

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		fmt.Println(err)
	}

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
		//dataMap[apgwName] = "id"
	}

	for key, value := range dataMap {
		fmt.Printf("APGW Name: %s, Firewall Policy ID: %s\n", key, value)
	}

}
