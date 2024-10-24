package main

import (
	"bufio"
	"context"
	"fmt"
	"gameserver/common/net/http"
	"gameserver/conf"
	"gameserver/core/database"
	"gameserver/core/logger"
	"gameserver/core/redis"
	"gameserver/excel"
	"gameserver/excel/cfgmgr"
	"gameserver/models/mmgr"
	"gameserver/models/models"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// 监听关闭信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel() // 取消上下文
	}()

	InitServer()

	cfgmgr.NewManager()
	cfg := excel.LevelConfig[1]
	logger.Logger.Info(fmt.Sprintf("level config:%T %v", cfg, cfg))

	go http.NewGinServer() // 将上下文传递给httpserver
	go monitor(ctx)        // 将上下文传递给monitor
	go scanner()

	logger.Logger.Info("appserver start success")

	// 等待上下文被取消
	<-ctx.Done()

	logger.Logger.Info("appserver is shutting down...")

	// 进行必要的清理工作
	mmgr.Stop()

	// 确保所有 goroutine 都已退出
	time.Sleep(3 * time.Second)

	logger.Logger.Info("appserver has shut down successfully")
}

func InitServer() {
	conf.InitConfig("./conf")
	logger.Init()
	redis.InitRedis()
	database.InitDatabase()
	models.InitAutoMigrate()
	mmgr.Start()
}

func scanner() {
	// 从标准输入流中接收输入数据
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "excel up" {
			cfg := excel.LevelConfig[1]
			logger.Logger.Info(fmt.Sprintf("level config:%T %v", cfg, cfg))
			continue
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
