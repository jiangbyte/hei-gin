// Package main æ˜¯æ•°æ®åº“è¿ç§»å…¥å£ï¼šåŸºäºŽ goose æ‰§è¡Œ up/down/statusã€‚
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

// è¿›ç¨‹å…¥å£ã€‚
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

// å‘æ ‡å‡†é”™è¯¯è¾“å‡ºåŽé€€å‡ºã€‚
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
	os.Exit(1)
}
