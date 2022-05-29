package week3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

// 定义一个简单的实现了http.Handler的结构体
type handler struct {
	content string
}

func NewHandler(context string) *handler {
	return &handler{
		content: context,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(h.content)
}

func TestApp(t *testing.T) {
	var srvs = []*Server{
		{
			name: "server1",
			srv: &http.Server{
				Addr:    ":8084",
				Handler: NewHandler("this is server1"),
			},
		},
		{
			name: "server2",
			srv: &http.Server{
				Addr:    ":8085",
				Handler: NewHandler("this is server2"),
			},
		},
		{
			name: "server3",
			srv: &http.Server{
				Addr:    ":8086", // 比如这里改为":8085", 就会产生一个error, 其他的两个server也会关闭
				Handler: NewHandler("this is server3"),
			},
		},
	}
	// 注册信号
	var sigs = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	ctx, cancel := context.WithCancel(context.Background())
	var opts = []Option{
		WithCancel(ctx, cancel),
		WithTimeout(time.Second),
		WithSrvs(srvs),
		WithSigs(sigs),
	}

	app := new(App)
	for _, opt := range opts {
		opt(app) // 配置App
	}

	// 启动
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
