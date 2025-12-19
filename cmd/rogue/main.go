package main

import (
	"flag"
	"log"
	"os"

	"rogue/internal/service"
	"rogue/internal/domain"
	"rogue/internal/adapters/storage"
	ui "rogue/internal/adapters/ui"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.LUTC)

	var (
		seed      = flag.Int64("seed", 0, "RNG seed (0 uses current time)")
		storePath = flag.String("store", "rogue_save.json", "path to JSON save/leaderboard store")
	)
	flag.Parse()

	st := storage.NewJSON(*storePath)
	gameApp := service.New(st, *seed, domain.Config{})

	if err := ui.Run(gameApp, ui.Options{Logger: logger}); err != nil {
		logger.Printf("fatal: %v", err)
		os.Exit(1)
	}
}
