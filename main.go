package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 自定义 handler
type helloHandler struct{}

// 实现 ServerHTTP 方法
func (handler *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Millisecond * 1000)
	w.Write([]byte("hello world1"))
}

type aboutHandler struct{}

// 实现 ServerHTTP 方法
func (handler *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(" About ...."))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(" welcome ...."))
}

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("hello world"))
	// })

	// 创建 web server
	// param1: 网络地址
	// param2：handler， 如果是nil，那么就是 DefaultServerMux
	// http.ListenAndServe("localhost:8080", nil)

	// 创建 web server
	mh := helloHandler{}
	a := aboutHandler{}
	server := http.Server{
		Addr: "localhost:8080",
		// Handler: &mh,
		Handler: nil,
		// Handler: http.NotFoundHandler(),
		// Handler: http.FileServer(http.Dir("wwwroot")),

	}
	// http://localhost:8080/hello
	http.Handle("/hello", &mh)
	http.Handle("/about", &a)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		// 返回切片
		// userAgent := header["User-Agent"]
		// 获取header字符串
		userAgentStr := header.Get("User-Agent")
		w.Write([]byte("Home... userAgent : " + userAgentStr))
	})
	// http://localhost:8080/welcome
	http.Handle("/welcome", http.HandlerFunc(welcome))
	// 几个常用的 handler
	// 请求响应 404
	// http://localhost:8080/404
	http.Handle("/404", http.NotFoundHandler())
	// 重定向
	// http://localhost:8080/redir
	http.Handle("/redir", http.RedirectHandler("http://localhost:8080/hello", 302))
	// 去掉指定的前缀，然后跳转到另一个handler
	// http://localhost:8080/ab
	http.Handle("/ab", http.StripPrefix("/ab", &mh))
	// 在指定时间内运行传入的handler，如果超时则返回 message
	// http://localhost:8080/timeout
	http.Handle("/timeout", http.TimeoutHandler(&mh, time.Millisecond*1500, "timeout error"))
	// 使用一个基于root文件系统响应请求
	// http://localhost:8080/index.html
	http.Handle("/", http.FileServer(http.Dir("wwwroot")))

	// 测试 query 参数
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		// 返回切片
		id := query["id"]
		fmt.Println(id)
		// 返回字符串
		name := query.Get("name")
		fmt.Println("name " + name)
	})

	// 测试 body
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		fmt.Fprintln(w, string(body))
	})

	// 测试 表单
	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		// 直接返回指定表单字段
		name := r.FormValue("name")
		fmt.Println(name)
		// 先解析request
		r.ParseForm()
		// 是个 map
		forms := r.Form
		// PostForm 中的字段只会获取表单中的数据， url中的数据不会获取到
		postForms := r.PostForm
		fmt.Println(forms)
		fmt.Println(postForms)
	})

	// multipartForm
	http.HandleFunc("/multipartForm", func(w http.ResponseWriter, r *http.Request) {
		// 直接返回指定表单字段
		name := r.PostFormValue("name")
		fmt.Println(name)
		// 先解析request
		r.ParseMultipartForm(1024)
		// 是个 struct, struct 内有两个map，第一个是表单数据，第二个是上传的文件
		forms := r.MultipartForm
		fmt.Fprintln(w, forms)

		// 接收上传文件，这里只接收第一个文件
		fileHeader := r.MultipartForm.File["upload"][0]
		file, err := fileHeader.Open()

		// 一个简便写法 接收上传文件，这里只接收第一个文件
		// file, _, err := r.FormFile("upload")

		if err == nil {
			data, err := ioutil.ReadAll(file)
			if err == nil {
				fmt.Fprintln(w, string(data))
			}
		}

	})

	server.ListenAndServe()

}
