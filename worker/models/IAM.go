package models

/*
그룹 쿼리
*/
// steampipe 쿼리 결과
type AdGroupResponse struct {
	Rows []AdGroup `json:"rows"`
}

// Rows 구조
type AdGroup struct {
	DisplayName        string `json:"display_name" gorm:"charset:utf8mb4;collate:utf8mb4_unicode_ci;"`
	Id                 string `json:"id" gorm:"primaryKey"`
	Description        string `json:"description" gorm:"charset:utf8mb4;collate:utf8mb4_unicode_ci;"`
	CreateDateTime     string `json:"create_datetime" gorm:"charset:utf8mb4;collate:utf8mb4_unicode_ci;"`
	ExpirationDateTime string `json:"expiration_datetime" gorm:"charset:utf8mb4;collate:utf8mb4_unicode_ci;"`
}

/*
dir Role쿼리
*/
//steampipe 쿼리 결과
type DirectoryRoleResponse struct {
	Rows []DirectoryRole `json:"rows"`
}

// Rows 구조
type DirectoryRole struct {
	Display_name string `json:"display_name" gorm:"charset:utf8mb4;collate:utf8mb4_unicode_ci;"`
	Id           string `json:"id" gorm:"primaryKey"`
}

/*
유저 쿼리
*/
//steampipe 쿼리 결과
type AzureADUserResponse struct {
	Rows     []AzureADUser `json:"rows"`
	Metadata *struct{}     `json:"metadata"`
}

// Rows 구조
type AzureADUser struct {
	DisplayName string `json:"display_name"`
	Mail        string `json:"mail"`
	Department  string `json:"department"`
	//MemberOf    string `json:"member_of"`
	MemberOf []GroupReference `json:"member_of"`
}

type GroupReference struct {
	OdataType string `json:"@odata.type"`
	ID        string `json:"id"`
}

type AzureADUserResult struct {
	Rows []AzureADUserReplace `json:"rows"`
}
type AzureADUserReplace struct {
	DisplayName string   `json:"display_name"`
	Mail        string   `json:"mail"`
	Department  string   `json:"department"`
	MemberOf    []string `json:"member_of"`
}
type ConvertMember struct {
	stringMemeber string `json:"stringMemeber"`
}

/*전송 데이터*/
// 유저 목록 grafana로 전송할 형태
type AzureADuserRequest struct {
	Rows []AzureADUserRemake `json:"rows"`
}

// Rows 구조
type AzureADUserRemake struct {
	DisplayName string `json:"display_name"`
	Mail        string `json:"mail"`
	Department  string `json:"department"`
	//MemberOf    string `json:"member_of"`
	MemberOf []GroupName `json:"member_of"`
}

type GroupName struct {
	ID        string `json:"id"`
	GroupName string `json:"group_name"`
}
