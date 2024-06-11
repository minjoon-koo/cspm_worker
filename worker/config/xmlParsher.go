package config

import (
	"encoding/xml"
	"io/ioutil"
	"worker/models"
)

func IamXmlParsher(queryTag string) string {
	// XML 파일 읽기
	data, err := ioutil.ReadFile("steamSQL/IAM.XML")
	if err != nil {
		//fmt.Println("Error reading file:", err)
		return "Error: reading file"
	}
	// XML 파싱
	var queries models.Queries
	err = xml.Unmarshal(data, &queries)
	if err != nil {
		//fmt.Println("Error parsing XML:", err)
		return "Error: parsing XML"
	}

	queryName := queryTag // 사용자에게 필요한 쿼리 이름
	var selectedQuery string
	for _, q := range queries.Queries {
		if q.Name == queryName {
			selectedQuery = q.Text
			break
		}
	}

	if selectedQuery == "" {
		//fmt.Println("Query not found")
		return "Error: Query not found"
	}

	return selectedQuery
}

func ApgwXmlParsher(queryTag string) string {
	// XML 파일 읽기
	data, err := ioutil.ReadFile("steamSQL/NET.XML")
	if err != nil {
		//fmt.Println("Error reading file:", err)
		return "Error: reading file"
	}
	// XML 파싱
	var queries models.Queries
	err = xml.Unmarshal(data, &queries)
	if err != nil {
		//fmt.Println("Error parsing XML:", err)
		return "Error: parsing XML"
	}

	queryName := queryTag // 사용자에게 필요한 쿼리 이름
	var selectedQuery string
	for _, q := range queries.Queries {
		if q.Name == queryName {
			selectedQuery = q.Text
			break
		}
	}

	if selectedQuery == "" {
		//fmt.Println("Query not found")
		return "Error: Query not found"
	}

	return selectedQuery
}
