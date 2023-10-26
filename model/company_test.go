package model

import "testing"

// 测试的go文件 _test 结尾，生产编译不会编译 _test 结尾的文件
//  方法 Test 开头， 参数 t *testing.T
func TestGetCompanyType(t *testing.T) {
	c := Company{
		ID:      1,
		Name:    "ABCD .LTD",
		Country: "China",
	}
	companyType := c.GetCompanyType()
	if companyType != "Limited Liability Company" {
		t.Errorf("Company's GetCompanyTYpe Method error")
	}
}
