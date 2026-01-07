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
	showStatistics  bool
	showStartMenu   bool
	rng             domain.RandomGenerator
}

func NewGame(ui UI) *Game {
	gameStateStore := storage.NewGameStateStorage("saves")
	hasSave := gameStateStore.HasSave()

	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)

	game := &Game{
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: storage.NewStatisticsStorage("saves"),
		gameStateStore:  gameStateStore,
		showStatistics:  false,
		showStartMenu:   hasSave,
		rng:             rng,
	}

	if !hasSave {
		world := domain.NewWorld(domain.Width, domain.Height)
		NewGenerator().Generate(world)
		game.World = world
		game.UpdateVisibility()
	}

	return game
}

func NewGameWithWorld(ui UI, world *domain.World) *Game {
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)

	game := &Game{
		World:           world,
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: storage.NewStatisticsStorage("saves"),
		gameStateStore:  storage.NewGameStateStorage("saves"),
		showStatistics:  false,
		showStartMenu:   false,
		rng:             rng,
	}

	game.UpdateVisibility()

	return game
}

func (g *Game) UpdateVisibility() {
	g.FOV.Update(g.World.Level, g.World.Player.Point)
}

func (g *Game) Run() {
	for g.isRunning {
		if g.showStartMenu {
			g.ui.DrawStartMenu(g.HasSaveGame())
		} else if g.showStatistics {
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
