package main

import (
	"flag"
	"log"
	"os"

	"rogue/internal/app"
	"rogue/internal/game"
	"rogue/internal/storage"
	tui "rogue/internal/ui"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.LUTC)

	var (
		seed      = flag.Int64("seed", 0, "RNG seed (0 uses current time)")
		storePath = flag.String("store", "rogue_save.json", "path to JSON save/leaderboard store")
	)
	flag.Parse()

	st := storage.NewJSON(*storePath)
	gameApp := app.New(st, *seed, game.Config{})

	if err := tui.Run(gameApp, tui.Options{Logger: logger}); err != nil {
		logger.Printf("fatal: %v", err)
		os.Exit(1)
	}
}
