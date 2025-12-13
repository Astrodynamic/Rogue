package main

import (
	"flag"
	"log"
	"os"

	"github.com/Astrodynamic/Rogue/internal/adapters/storage"
	"github.com/Astrodynamic/Rogue/internal/adapters/tui"
	"github.com/Astrodynamic/Rogue/internal/app/service"
	"github.com/Astrodynamic/Rogue/internal/app/usecase"
	"github.com/Astrodynamic/Rogue/internal/domain/game"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.LUTC)

	var (
		seed      = flag.Int64("seed", 0, "RNG seed (0 uses current time)")
		storePath = flag.String("store", "rogue_save.json", "path to JSON save/leaderboard store")
	)
	flag.Parse()

	st := storage.NewJSON(*storePath)
	svc := service.New(st, st)
	uc := usecase.New(svc, *seed, game.Config{})

	if err := tui.Run(uc, tui.Options{Logger: logger}); err != nil {
		logger.Printf("fatal: %v", err)
		os.Exit(1)
	}
}
