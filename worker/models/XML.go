package models

import "encoding/xml"

type Query struct {
	Name string `xml:"name,attr"`
	Text string `xml:",chardata"`
}

type Queries struct {
	XMLName xml.Name `xml:"Queries"`
	Queries []Query  `xml:"Query"`
}
