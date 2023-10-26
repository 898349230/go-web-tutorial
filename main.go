package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"github.com/898349230/go-web-tutorial/middleware"
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

func writeExample(w http.ResponseWriter, r *http.Request) {
	str := `<html> 
	<head><title>Go Web</title></head>
	<body><h1>Hello World</h1></body>
	</html>`
	w.Write([]byte(str))
}

func writeHeader(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	str := `<html> 
	<head><title>5501</title></head>
	<body><h1>Hello World</h1></body>
	</html>`
	w.Write([]byte(str))
}

func redirectHand(w http.ResponseWriter, r *http.Request) {
	// 需要再设置 WriteHeader 前设置 Location，调用完 WriteHeader 后无法设置header
	w.Header().Set("Location", "https://www.bilibili.com/")
	w.WriteHeader(302)
}
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	post := &Post{
		User:  "张三",
		Title: "心情号",
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}

type Post struct {
	User  string
	Title string
}

type Comapny struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
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
		// 设置中间件
		// Handler: new(middleware.AuthMiddleware),
		// 设置嵌套中间件
		Handler: &middleware.AuthMiddleware{
			Next: new(middleware.TimeoutMiddleware),
		},
		// Handler: nil,
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
		// sleep 3秒，测试超时中间件
		time.Sleep(3 * time.Second)
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

	// http://localhost:8080/writeHeader
	http.HandleFunc("/write", writeExample)
	// http://localhost:8080/writeHeader
	http.HandleFunc("/writeHeader", writeHeader)
	// http://localhost:8080/redirect
	http.HandleFunc("/redirect", redirectHand)
	// http://localhost:8080/json
	http.HandleFunc("/json", jsonHandler)

	// 内置的 Response
	// NotFound 函数，包装一个404状态码和一个额外的信息
	// ServeFile 函数，从文件系统提供文件，返回给请求者
	// ServeContent 函数，可以把io.ReadSeeker接口的任何东西里面的内容返回给请求者
	// Redirect 函数，告诉客户端重定向到另一个URL

	// 解析模版
	// http://localhost:8080/process
	http.HandleFunc("/process", process)
	// http://localhost:8080/process2
	http.HandleFunc("/process2", process2)
	// http://localhost:8080/process3
	http.HandleFunc("/process3", process3)
	// http://localhost:8080/processMethod
	http.HandleFunc("/processMethod", processMethod)

	// 接收 返回 json
	http.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// json解码器
			dec := json.NewDecoder(r.Body)
			company := Comapny{}
			err := dec.Decode(&company)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println(company)
			// json 编码器
			enc := json.NewEncoder(w)
			// 将接收到的数据返回去
			err = enc.Encode(company)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	server.ListenAndServe()
	// http.ListenAndServe("localhost:8080", new(middleware.AuthMiddleware))

}

// 解析模板
func process(w http.ResponseWriter, r *http.Request) {
	// t3, _ := template.ParseGlob("/*")

	// 解析模板
	t, _ := template.ParseFiles("templ.html")
	// 执行单个模版
	t.Execute(w, "Hello World haha ")

	// 解析多个html
	ts, _ := template.ParseFiles("t1.html", "t2.html")
	// 执行指定的html
	ts.ExecuteTemplate(w, "h2.html", "Hello World")

	// 查找指定文件
	// file := ts.Lookup("h1.html")

}

// 解析模板
func process2(w http.ResponseWriter, r *http.Request) {

	// 解析模板
	t, _ := template.ParseFiles("templ2.html")
	// 执行单个模版
	rand.Seed(time.Now().Unix())
	t.Execute(w, rand.Intn(10) > 5)
}

// 遍历
func process3(w http.ResponseWriter, r *http.Request) {

	// 解析模板
	t, _ := template.ParseFiles("templ3.html")
	dayOfWeek := []string{"Mon", "Tue", "Wed"}
	t.Execute(w, dayOfWeek)
}

// 模版自定义函数
func processMethod(w http.ResponseWriter, r *http.Request) {
	// 模版函数名称 fdate
	funcMap := template.FuncMap{"fdate": formateDate}
	t := template.New("templ4.html").Funcs(funcMap)
	t.ParseFiles("templ4.html")
	t.Execute(w, time.Now())
}

// 时间格式化
func formateDate(t time.Time) string {
	layout := "2006-01-02"
	return t.Format(layout)
}
