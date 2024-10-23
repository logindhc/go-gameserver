package main

import (
	"bufio"
	"context"
	"fmt"
	"gameserver/common/logger"
	"gameserver/common/net/http"
	"gameserver/common/persistence"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	gomaxprocs := runtime.GOMAXPROCS(runtime.NumCPU())
	logger.Logger.Info(fmt.Sprintf("appserver start gomaxprocs %v", gomaxprocs))
	ctx, cancel := context.WithCancel(context.Background())
	// 监听关闭信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Logger.Info("Received shutdown signal, shutting down...")
		cancel() // 取消上下文
	}()

	go http.NewGinServer() // 将上下文传递给httpserver
	go monitor(ctx)        // 将上下文传递给monitor

	logger.Logger.Info("appserver start success")

	// 等待上下文被取消
	<-ctx.Done()
	logger.Logger.Info("appserver is shutting down...")

	// 进行必要的清理工作
	persistence.Stop()

	// 确保所有 goroutine 都已退出
	time.Sleep(3 * time.Second)

	logger.Logger.Info("appserver has shut down successfully")
}

func scanner() {
	// 从标准输入流中接收输入数据
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "excel up" {

			break
		}
	}
}

func monitor(ctx context.Context) {
	t := time.NewTicker(60 * time.Second)
	defer t.Stop() // 确保 ticker 资源被释放

	for {
		select {
		case <-ctx.Done(): // 检查上下文是否被取消
			return
		case <-t.C:
			logger.Logger.Info(fmt.Sprintf("monitor goroutine count: %d\n", runtime.NumGoroutine()))
		}
	}
}
