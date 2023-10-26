package model

import (
	"encoding/json"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/company", handleCompany)
}

// 待测试方法
func handleCompany(w http.ResponseWriter, r *http.Request) {
	c := Company{
		ID:      2,
		Name:    "Google",
		Country: "US",
	}

	// json 编码写到 writer中
	enc := json.NewEncoder(w)
	enc.Encode(c)
}
