package service

import (
	"math/rand/v2"
	"time"

	"rogue/internal/adapters/storage"
	"rogue/internal/domain"
)

type Game struct {
	World           *domain.World
	FOV             *FOV
	isRunning       bool
	ui              UI
	selection       *SelectionModel
	statisticsStore *storage.StatisticsStorage
	gameStateStore  *storage.GameStateStorage
	rng             domain.RandomGenerator
	generator       *Generator
}

func NewGame(ui UI) *Game {
	gameStateStore := storage.NewGameStateStorage("saves")

	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)
	generator := NewGeneratorWithRNG(rng)

	game := &Game{
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: storage.NewStatisticsStorage("saves"),
		gameStateStore:  gameStateStore,
		rng:             rng,
		generator:       generator,
	}

	return game
}

func NewGameWithWorld(ui UI, world *domain.World) *Game {
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)
	generator := NewGeneratorWithRNG(rng)

	game := &Game{
		World:           world,
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: storage.NewStatisticsStorage("saves"),
		gameStateStore:  storage.NewGameStateStorage("saves"),
		rng:             rng,
		generator:       generator,
	}

	game.UpdateVisibility()

	return game
}

func (g *Game) UpdateVisibility() {
	g.FOV.Update(g.World.Level, g.World.Player.Point)
}

func (g *Game) Run() {
	for g.isRunning {
		if g.ui.IsStartMenu() {
			g.ui.DrawStartMenu(g.HasSaveGame(), g.ui.GetMenuOption())
		} else if g.ui.IsStatistics() {
			playthroughs, _ := g.statisticsStore.GetLeaderboard(0)
			g.ui.DrawStatistics(playthroughs)
		} else {
			selection := SelectionState{
				Model: g.selection,
			}
			g.ui.Draw(g.World, selection)
		}
		g.handle(g.ui.Input())
	}
}
