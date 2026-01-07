package service

import (
	"math/rand/v2"
	"time"

	"rogue/internal/domain"
)

type Game struct {
	World     *domain.World
	FOV       *FOV
	isRunning bool
	ui        UI
	sel       *SelectionModel
	stats     StatsStore
	state     StateStore
	rng       domain.RandomGenerator
	gen       *Generator
}

func NewGame(ui UI, state StateStore, stats StatsStore) *Game {
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)
	gen := NewGenWithRNG(rng)

	game := &Game{
		FOV:       NewFOV(),
		isRunning: true,
		ui:        ui,
		stats:     stats,
		state:     state,
		rng:       rng,
		gen:       gen,
	}

	return game
}

func NewGameWithW(ui UI, world *domain.World, state StateStore, stats StatsStore) *Game {
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)
	gen := NewGenWithRNG(rng)

	game := &Game{
		World:     world,
		FOV:       NewFOV(),
		isRunning: true,
		ui:        ui,
		stats:     stats,
		state:     state,
		rng:       rng,
		gen:       gen,
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
			playthroughs, _ := g.stats.GetLeaderboard(0)
			g.ui.DrawStatistics(playthroughs)
		} else if g.World != nil {
			sel := SelectionState{
				Model: g.sel,
			}
			g.ui.Draw(g.World, sel)
		}
		g.handle(g.ui.Input())
	}
}
