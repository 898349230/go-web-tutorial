package middleware

import (
	"context"
	"net/http"
	"time"
)

// 超时中间件
type TimeoutMiddleware struct {
	Next http.Handler
}

// 实现 Handler 接口
func (tm *TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}
	// 获取当前 Context
	ctx := r.Context()
	// 修改当前 Context， 3秒超时
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	// 将修改的 Context 替换原来的 Context
	r.WithContext(ctx)
	// 空 struct 的 channel
	ch := make(chan struct{})
	go func() {
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		return
	case <-ctx.Done():
		// 超时
		w.WriteHeader(http.StatusRequestTimeout)
	}

}
