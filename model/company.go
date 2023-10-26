package model

import "strings"

type Company struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

func (c *Company) GetCompanyType() (result string) {
	if strings.HasSuffix(c.Name, ".LTD") {
		result = "Limited Liability Company"
	} else {
		result = "Others"
	}
	return
}
