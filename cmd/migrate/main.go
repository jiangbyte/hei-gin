package main

import (
	"flag"
	"fmt"
	"log"

	plugin_client "hei-gin/plugins/plugin-client"
	plugin_im "hei-gin/plugins/plugin-im"
	plugin_sys "hei-gin/plugins/plugin-sys"
	"hei-gin/sdk/config"
	"hei-gin/sdk/infra/db"
)

func main() {
	plugin_sys.RegisterMigrations()
	plugin_client.RegisterMigrations()
	plugin_im.RegisterMigrations()

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

	db.Freeze()

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
