package service

import (
	"math/rand/v2"
	"time"

	"rogue/internal/domain"
)

type Game struct {
	World           *domain.World
	FOV             *FOV
	isRunning       bool
	ui              UI
	selection       *SelectionModel
	statisticsStore StatsStore
	gameStateStore  StateStore
	rng             domain.RandomGenerator
	generator       *Generator
}

func NewGame(ui UI, gameStateStore StateStore, statisticsStore StatsStore) *Game {
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)
	generator := NewGenWithRNG(rng)

	game := &Game{
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: statisticsStore,
		gameStateStore:  gameStateStore,
		rng:             rng,
		generator:       generator,
	}

	return game
}

func NewGameWithW(ui UI, world *domain.World, gameStateStore StateStore, statisticsStore StatsStore) *Game {
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)
	generator := NewGenWithRNG(rng)

	game := &Game{
		World:           world,
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: statisticsStore,
		gameStateStore:  gameStateStore,
		rng:             rng,
		generator:       generator,
	}

	game.UpdateVisibility()

	return game
}

func (g *Game) UpdateVisibility() {
	g.FOV.Update(g.World.Level, g.World.Player.Point)
}

func (g *Game) checkPlayerHealth() {
	if g.World != nil && g.World.Player.Health <= 0 {
		g.handlePlayerDeath()
	}
}

func (g *Game) Run() {
	for g.isRunning {
		if g.ui.IsStartMenu() {
			g.ui.DrawStartMenu(g.HasSaveGame(), g.ui.GetMenuOption())
		} else if g.ui.IsNameInput() {
			var emptyWorld *domain.World
			g.ui.Draw(emptyWorld, SelectionState{})
		} else if g.ui.IsStatistics() {
			playthroughs, _ := g.statisticsStore.GetLeaderboard(0)
			g.ui.DrawStatistics(playthroughs)
		} else if g.World != nil {
			selection := SelectionState{
				Model: g.selection,
			}
			g.ui.Draw(g.World, selection)
		}
		g.handle(g.ui.Input())
	}
}
