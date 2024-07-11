package main

import (
	query "worker/config"
	"worker/controller/NET"
)

var sql string = query.ApgwXmlParsher("getFirewallPolicy")
var url string = "https://portal.azure.com/#@x2soft.onmicrosoft.com/resource"

func main() {
	//db.Connection()
	NET.ApgwFrontendInfo()
	//controller.GetAppGatewayInfo()

	//NET.UpdateFrontendResource()
}
