// Package main æ˜¯å•ä½“ API è¿›ç¨‹å…¥å£ï¼šåŠ è½½é…ç½®ã€è£…é…æ¨¡å—ã€HTTP ä¸Žè¿›ç¨‹å†…å®šæ—¶ä»»åŠ¡ã€‚
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hei-gin/internal/app"
	_ "hei-gin/internal/modules/all"
)

// è¿›ç¨‹å…¥å£ã€‚
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

	api := app.NewAPI(d)
	ctx := context.Background()
	if err := api.Start(ctx); err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = api.Stop(shutdown)
}
