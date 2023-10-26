package middleware

import "net/http"

// 权限验证中间件
type AuthMiddleware struct {
	Next http.Handler
}

// 实现 Handler 接口
func (am *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if am.Next == nil {
		am.Next = http.DefaultServeMux
	}

	auth := r.Header.Get("Authorization")
	if auth != "" {
		am.Next.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
