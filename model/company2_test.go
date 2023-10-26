package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试的go文件 _test 结尾，生产编译不会编译 _test 结尾的文件
//  方法 Test 开头， 参数 t *testing.T
func TestHandleCompany(t *testing.T) {
	// 模拟请求
	r := httptest.NewRequest(http.MethodGet, "/company", nil)
	// 模拟响应
	w := httptest.NewRecorder()

	// 发起请求
	handleCompany(w, r)

	// 获取 body
	result, _ := ioutil.ReadAll(w.Result().Body)

	c := Company{}
	json.Unmarshal(result, &c)

	if c.ID != 2 {
		t.Errorf("Fail to handle company correctly!")
	}

}
