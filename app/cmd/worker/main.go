// Package main 是 Worker 进程入口：加载配置、装配模块并跑定时任务。
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hei-gin/app/internal/app"
	_ "hei-gin/app/internal/modules/all"
)

// 进程入口。
func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	cfg := app.LoadOrDie(cfgPath)
	d, err := app.OpenInfra(cfg)
	if err != nil {
		panic(err)
	}
	app.AttachRegisteredModules(d)

	w := app.NewWorker(d)
	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = w.Stop(shutdown)
}
