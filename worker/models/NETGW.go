package models

/*
1차 steampipe 검색 결과
*/
type SteamSQLappGateway struct {
	Rows []AppGateways `json:"rows"`
}

type AppGateways struct {
}

type AppGatewayHealth struct {
	ProbeName   string `json:"name"`
	ProbeStatus string `json:"status"`
}

/*
상태
*/
type BackendPollResponse struct {
	Rows []BackendPool `json:"rows"`
}
type BackendPool struct {
	AppGatewayName string `json:"apgw_name"`
	BackendName    string `json:"backend_name"`
	Ipaddress      string `json:"ipaddress"`
	BackendID      string `json:"backend_id" gorm:"primaryKey"`
}

/*
Azure az cli
*/

type StatuesResponse struct {
	Address string `json:"address"`
	Health  string `json:"health"`
}

type BackendStatus struct {
	ApgwName string            `json:"apgw_name"`
	Statues  []StatuesResponse `json:"statues"`
}

/*
방화벽 정보
*/
type WafInformation struct {
	ApgwName           string `json:"apgw_name"`
	FirewallPolicyName string `json:"firewall_policy_name"`
	FirewallResources  string `json:"firewall_resources"`
	FirewallPolicyID   string `json:"firewall_policy_id" gorm:"primaryKey"`
}
type WafInformations struct {
	WafInfo []WafInformation `json:"waf_info"`
}
