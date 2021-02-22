package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func testRequest(t *testing.T, handler *web.ControllerRegister, requestIP, method, path string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	r.Header.Set("X-Real-Ip", requestIP)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", requestIP, method, path, w.Code, code)
	}
}

func TestLimiter(t *testing.T) {
	opts := []LimiterOption{WithRate(1 * time.Millisecond), WithCapacity(1)}
	handler := web.NewControllerRegister()
	handler.InsertFilter("/foo/*", web.BeforeRouter, NewLimiter(opts))

	handler.Any("*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "127.0.0.1", "GET", "/foo/1", 200)
	testRequest(t, handler, "127.0.0.1", "GET", "/foo/1", 429)
	testRequest(t, handler, "127.0.0.2", "GET", "/foo/1", 200)
	time.Sleep(1 * time.Millisecond)
	testRequest(t, handler, "127.0.0.1", "GET", "/foo/1", 200)
}

func Benchmark_WithoutLimiter(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := web.NewControllerRegister()
	web.BConfig.RunMode = web.PROD
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			handler.ServeHTTP(recorder, r)
		}
	})
}

func Benchmark_WithLimiter(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := web.NewControllerRegister()
	web.BConfig.RunMode = web.PROD
	opts := []LimiterOption{WithRate(1 * time.Millisecond), WithCapacity(100)}
	handler.InsertFilter("*", web.BeforeRouter, NewLimiter(opts))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			handler.ServeHTTP(recorder, r)
		}
	})
}
