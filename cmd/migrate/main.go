// Package main 是数据库迁移入口：基于 goose 执行 up/down/status。
//
// Author: Charlie
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"hei-gin/internal/framework/core/config"
)

// 进程入口。
func main() {
	cfgPath := flag.String("config", "config.yaml", "config file")
	dir := flag.String("dir", "migrations", "migrations directory")
	flag.Parse()
	args := flag.Args()
	cmd := "up"
	if len(args) > 0 {
		cmd = args[0]
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fatal(err)
	}
	db, err := sql.Open("pgx", cfg.DB.URL)
	if err != nil {
		fatal(err)
	}
	defer db.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		fatal(err)
	}
	switch cmd {
	case "up":
		err = goose.Up(db, *dir)
	case "down":
		err = goose.Down(db, *dir)
	case "status":
		err = goose.Status(db, *dir)
	default:
		fatal(fmt.Errorf("unknown command %q (up|down|status)", cmd))
	}
	if err != nil {
		fatal(err)
	}
}

// 向标准错误输出后退出。
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
	os.Exit(1)
}
