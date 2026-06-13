package main

import (
	"flag"
	"fmt"
	"log"

	"hei-gin/sdk/config"
	"hei-gin/sdk/infra/db"

	// Blank-import plugins to trigger model + seed self-registration via init()
	_ "hei-gin/plugins/plugin-client/user"
	_ "hei-gin/plugins/plugin-im/broadcast"
	_ "hei-gin/plugins/plugin-im/friend"
	_ "hei-gin/plugins/plugin-im/group"
	_ "hei-gin/plugins/plugin-im/message"
	_ "hei-gin/plugins/plugin-im/model"
	_ "hei-gin/plugins/plugin-sys/banner"
	_ "hei-gin/plugins/plugin-sys/config"
	_ "hei-gin/plugins/plugin-sys/dict"
	_ "hei-gin/plugins/plugin-sys/file"
	_ "hei-gin/plugins/plugin-sys/group"
	_ "hei-gin/plugins/plugin-sys/home"
	_ "hei-gin/plugins/plugin-sys/log"
	_ "hei-gin/plugins/plugin-sys/notice"
	_ "hei-gin/plugins/plugin-sys/org"
	_ "hei-gin/plugins/plugin-sys/position"
	_ "hei-gin/plugins/plugin-sys/resource"
	_ "hei-gin/plugins/plugin-sys/role"
	_ "hei-gin/plugins/plugin-sys/user"
)

func main() {
	skipSeed := flag.Bool("skip-seed", false, "skip seeding initial data")
	flag.Parse()

	if err := config.FindAndLoad(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := config.C.ValidateMigration(); err != nil {
		log.Fatalf("invalid migration config: %v", err)
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	defer db.Close()

	models := db.GetModels()
	if len(models) == 0 {
		fmt.Println("No models registered, skipping migration")
		return
	}
	if err := db.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to apply migration: %v", err)
	}
	fmt.Println("✓ Migration applied successfully")

	if !*skipSeed {
		if err := db.RunSeeds(); err != nil {
			log.Fatalf("failed to run seeds: %v", err)
		}
		fmt.Println("✓ Seeds applied successfully")
	}
}
