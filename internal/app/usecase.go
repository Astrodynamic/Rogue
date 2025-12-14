package app

import (
	"errors"
	"math/rand"
	"time"
	"unicode"
	"unicode/utf8"

	"rogue/internal/game"
)

type InputKind uint8

const (
	InputNone InputKind = iota
	InputUp
	InputDown
	InputLeft
	InputRight
	InputAttack
	InputQuit
	InputSave
	InputInvWeapon
	InputInvFood
	InputInvElixir
	InputInvScroll
	InputViewStats
	InputViewLeaderboard
	InputHelp
	InputDigit0
	InputDigit1
	InputDigit2
	InputDigit3
	InputDigit4
	InputDigit5
	InputDigit6
	InputDigit7
	InputDigit8
	InputDigit9
	InputConfirmYes
	InputConfirmNo
	InputText
	InputBackspace
)

type Input struct {
	Kind InputKind
	Rune rune
}

type Mode uint8

const (
	ModePlay Mode = iota
	ModeSelectSlot
	ModeSlotAction
	ModeConfirmOverwrite
	ModeConfirmDelete
	ModeSelectWeapon
	ModeSelectFood
	ModeSelectElixir
	ModeSelectScroll
	ModeStats
	ModeLeaderboard
	ModeHelp
	ModeEnterName
	ModeGameOver
	ModeWin
)

type ViewModel struct {
	Title    string
	Help     string
	Mode     Mode
	Level    int
	MapW     int
	MapH     int
	Tiles    [][]game.Tile
	Explored [][]bool
	Visible  [][]bool

	PlayerPos game.Point
	Player    game.Character
	Weapon    *game.Item
	Backpack  game.Backpack
	Enemies   []game.Enemy
	Rooms     []game.Room

	Messages []string

	// For menus.
	MenuTitle string
	MenuItems []string

	StatsLines []string
	BoardLines []string
}

type App struct {
	storage Storage

	rng  *rand.Rand
	sess *game.GameSession

	mode Mode

	// Selected inventory kind.
	pending game.ActionKind

	name string

	cfg game.Config

	activeSlot int

	selectedSlot    int
	selectedHasSave bool
}

// SetWorldSize updates the default world size used when starting a NEW session.
// Loaded sessions keep their saved map size.
func (a *App) SetWorldSize(width, height int) {
	a.cfg = game.Config{Width: width, Height: height}.Normalize()
	// If we're still in slot selection (pre-run), keep the seed but regenerate the initial level for a nicer first impression.
	if a.mode == ModeSelectSlot && a.sess != nil {
		a.sess = game.NewSession(a.sess.Seed, a.cfg)
	}
}

// New creates a new game application instance.
func New(storage Storage, seed int64, cfg game.Config) *App {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))
	return &App{
		storage: storage,
		rng:     rng,
		sess:    game.NewSession(seed, cfg),
		mode:    ModeSelectSlot,
		cfg:     cfg.Normalize(),
	}
}

func (a *App) Handle(in Input) (ViewModel, bool, error) {
	if a.sess == nil {
		return ViewModel{}, true, errors.New("nil session")
	}

	switch a.mode {
	case ModeSelectSlot:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		if slot, ok := digitIndex(in); ok && slot >= 1 && slot <= 3 {
			a.selectedSlot = slot
			_, has, _ := a.storage.LoadSession(slot)
			a.selectedHasSave = has
			a.mode = ModeSlotAction
			return a.view(), false, nil
		}
		return a.view(), false, nil
	case ModeSlotAction:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		// 1=Continue/New, 2=New (overwrite), 3=Delete, 0=Back
		if d, ok := digitIndex(in); ok {
			switch d {
			case 0:
				a.mode = ModeSelectSlot
				return a.view(), false, nil
			case 1:
				if a.selectedHasSave {
					a.activeSlot = a.selectedSlot
					if saved, ok2, err := a.storage.LoadSession(a.activeSlot); err == nil && ok2 {
						a.sess = game.FromSnapshot(saved.Snapshot)
						a.rng = rand.New(rand.NewSource(a.sess.Seed))
						a.name = saved.Name
						if normalizeName(a.name) == "" {
							a.mode = ModeEnterName
						} else {
							a.mode = ModePlay
						}
						a.sess.Messages = append(a.sess.Messages, "Loaded slot "+itoa(a.activeSlot)+".")
						return a.view(), false, nil
					}
					a.mode = ModeEnterName
					return a.view(), false, nil
				}
				// Empty slot -> start new
				a.activeSlot = a.selectedSlot
				a.mode = ModeEnterName
				return a.view(), false, nil
			case 2:
				if a.selectedHasSave {
					a.mode = ModeConfirmOverwrite
					return a.view(), false, nil
				}
				a.activeSlot = a.selectedSlot
				a.mode = ModeEnterName
				return a.view(), false, nil
			case 3:
				if a.selectedHasSave {
					a.mode = ModeConfirmDelete
					return a.view(), false, nil
				}
				return a.view(), false, nil
			}
		}
		return a.view(), false, nil
	case ModeConfirmOverwrite:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		if in.Kind == InputConfirmNo {
			a.mode = ModeSlotAction
			return a.view(), false, nil
		}
		if in.Kind == InputConfirmYes {
			a.activeSlot = a.selectedSlot
			_ = a.storage.ClearSession(a.activeSlot)
			a.selectedHasSave = false
			a.mode = ModeEnterName
			return a.view(), false, nil
		}
		return a.view(), false, nil
	case ModeConfirmDelete:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		if in.Kind == InputConfirmNo {
			a.mode = ModeSlotAction
			return a.view(), false, nil
		}
		if in.Kind == InputConfirmYes {
			_ = a.storage.ClearSession(a.selectedSlot)
			a.selectedHasSave = false
			a.mode = ModeSelectSlot
			return a.view(), false, nil
		}
		return a.view(), false, nil
	case ModeStats:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		a.mode = ModePlay
		return a.view(), false, nil
	case ModeLeaderboard:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		a.mode = ModePlay
		return a.view(), false, nil
	case ModeHelp:
		if in.Kind == InputQuit {
			return a.view(), true, nil
		}
		a.mode = ModePlay
		return a.view(), false, nil
	case ModeEnterName:
		switch in.Kind {
		case InputQuit:
			return a.view(), true, nil
		case InputConfirmYes:
			a.name = normalizeName(a.name)
			if a.name == "" {
				a.name = "player"
			}
			if a.activeSlot == 0 {
				a.activeSlot = 1
			}
			seed := time.Now().UnixNano()
			a.rng = rand.New(rand.NewSource(seed))
			a.sess = game.NewSession(seed, a.cfg)
			a.mode = ModePlay
			return a.view(), false, nil
		case InputBackspace:
			a.name = backspace(a.name)
			return a.view(), false, nil
		case InputText:
			if len([]rune(a.name)) < 16 && isNameRune(in.Rune) {
				a.name += string(in.Rune)
			}
			return a.view(), false, nil
		default:
			return a.view(), false, nil
		}
	case ModeSelectWeapon, ModeSelectFood, ModeSelectElixir, ModeSelectScroll:
		idx, ok := digitIndex(in)
		if !ok {
			if in.Kind == InputQuit || in.Kind == InputConfirmNo {
				a.mode = ModePlay
			}
			return a.view(), in.Kind == InputQuit, nil
		}
		a.applySelection(idx)
		a.mode = ModePlay
		return a.view(), false, nil
	case ModeGameOver:
		// Any key restarts.
		a.startNewRun()
		return a.view(), false, nil
	case ModeWin:
		// Any key exits.
		return a.view(), true, nil
	}

	// Play mode.
	switch in.Kind {
	case InputQuit:
		a.saveNow(false)
		return a.view(), true, nil
	case InputSave:
		a.saveNow(true)
	case InputHelp:
		a.mode = ModeHelp
	case InputUp:
		a.step(game.Action{Kind: game.ActionMove, Dx: 0, Dy: -1})
	case InputDown:
		a.step(game.Action{Kind: game.ActionMove, Dx: 0, Dy: 1})
	case InputLeft:
		a.step(game.Action{Kind: game.ActionMove, Dx: -1, Dy: 0})
	case InputRight:
		a.step(game.Action{Kind: game.ActionMove, Dx: 1, Dy: 0})
	case InputAttack:
		a.step(game.Action{Kind: game.ActionAttack})
	case InputInvWeapon:
		a.mode = ModeSelectWeapon
		a.pending = game.ActionEquipWeapon
	case InputInvFood:
		a.mode = ModeSelectFood
		a.pending = game.ActionUseFood
	case InputInvElixir:
		a.mode = ModeSelectElixir
		a.pending = game.ActionUseElixir
	case InputInvScroll:
		a.mode = ModeSelectScroll
		a.pending = game.ActionUseScroll
	case InputViewStats:
		a.mode = ModeStats
	case InputViewLeaderboard:
		a.mode = ModeLeaderboard
	}

	return a.view(), false, nil
}

func (a *App) step(action game.Action) {
	prevLevel := a.sess.LevelDepth
	advanced, dead, won, _ := a.sess.Step(a.rng, action)
	if advanced && a.sess.LevelDepth != prevLevel {
		// Autosave only on level transitions (reaching the exit).
		a.saveNow(false)
	}
	if won {
		a.finishRun(true)
		a.mode = ModeWin
	}
	if dead {
		a.finishRun(false)
		a.mode = ModeGameOver
	}
}

func (a *App) saveNow(showMessage bool) {
	if a.storage == nil || a.sess == nil {
		return
	}
	name := normalizeName(a.name)
	if name == "" {
		// Don't create anonymous saves; name entry flow happens first anyway.
		return
	}
	if a.activeSlot == 0 {
		a.activeSlot = 1
	}
	_ = a.storage.SaveSession(a.activeSlot, SavedSession{Snapshot: a.sess.Snapshot(), Name: name})
	if showMessage {
		a.sess.Messages = append(a.sess.Messages, "Saved (slot "+itoa(a.activeSlot)+").")
	}
}

func (a *App) finishRun(won bool) {
	name := normalizeName(a.name)
	if name == "" {
		name = "player"
	}
	_ = a.storage.RecordRun(RunResult{
		Name:         name,
		DeepestLevel: a.sess.Stats.DeepestLevelReached,
		Treasure:     a.sess.Stats.TotalTreasure,
		Stats:        a.sess.Stats,
	})
	if won {
		a.sess.Messages = append(a.sess.Messages, "You won! Press any key to exit.")
	} else {
		a.sess.Messages = append(a.sess.Messages, "Game over. Press any key to restart.")
	}
}

func (a *App) startNewRun() {
	a.mode = ModeSelectSlot
	a.activeSlot = 0
}

func (a *App) applySelection(idx int) {
	switch a.pending {
	case game.ActionEquipWeapon:
		if idx == 0 {
			a.step(game.Action{Kind: game.ActionUnequipWeapon})
			return
		}
		a.step(game.Action{Kind: game.ActionEquipWeapon, Idx: idx - 1})
	case game.ActionUseFood:
		a.step(game.Action{Kind: game.ActionUseFood, Idx: idx - 1})
	case game.ActionUseElixir:
		a.step(game.Action{Kind: game.ActionUseElixir, Idx: idx - 1})
	case game.ActionUseScroll:
		a.step(game.Action{Kind: game.ActionUseScroll, Idx: idx - 1})
	}
}

func digitIndex(in Input) (int, bool) {
	switch in.Kind {
	case InputDigit0:
		return 0, true
	case InputDigit1:
		return 1, true
	case InputDigit2:
		return 2, true
	case InputDigit3:
		return 3, true
	case InputDigit4:
		return 4, true
	case InputDigit5:
		return 5, true
	case InputDigit6:
		return 6, true
	case InputDigit7:
		return 7, true
	case InputDigit8:
		return 8, true
	case InputDigit9:
		return 9, true
	default:
		return 0, false
	}
}

func (a *App) view() ViewModel {
	vis := a.sess.ComputeVisible()
	var wp *game.Item
	if wid := a.sess.Player.Equipment.Get(game.SlotWeapon); wid != 0 {
		if it, ok := a.sess.Backpack.FindByID(wid); ok && it.Type == game.ItemWeapon {
			cp := it
			wp = &cp
		}
	}
	mw := len(a.sess.Level.Tiles[0])
	mh := len(a.sess.Level.Tiles)
	vm := ViewModel{
		Title:     "Rogue",
		Help:      "Move WASD  Attack Space  Inv h/j/k/e  Save Ctrl+S  (auto on exit+quit)  Stats t  Board l  Help ?  Quit q",
		Mode:      a.mode,
		Level:     a.sess.LevelDepth,
		MapW:      mw,
		MapH:      mh,
		Tiles:     a.sess.Level.Tiles,
		Explored:  a.sess.Explored,
		Visible:   vis,
		PlayerPos: a.sess.PlayerPos,
		Player:    a.sess.Player,
		Weapon:    wp,
		Backpack:  a.sess.Backpack,
		Enemies:   append([]game.Enemy(nil), a.sess.Enemies...),
		Rooms:     a.sess.Level.Rooms,
		Messages:  append([]string(nil), a.sess.Messages...),
	}

	switch a.mode {
	case ModeSelectSlot:
		vm.MenuTitle = "Select save slot (1-3). Empty slot starts a new game."
		if slots, err := a.storage.ListSaveSlots(); err == nil {
			items := make([]string, 0, len(slots))
			for _, sl := range slots {
				if sl.Empty {
					items = append(items, "Slot "+itoa(sl.Slot)+": <empty>")
				} else {
					items = append(items, "Slot "+itoa(sl.Slot)+": "+sl.Name+" lvl="+itoa(sl.Level)+" gold="+itoa(sl.Treasure))
				}
			}
			vm.MenuItems = items
		} else {
			vm.MenuItems = []string{"Slot 1", "Slot 2", "Slot 3"}
		}
	case ModeSlotAction:
		if a.selectedHasSave {
			vm.MenuTitle = "Slot " + itoa(a.selectedSlot) + ": 1=Continue  2=New(overwrite)  3=Delete  0=Back"
		} else {
			vm.MenuTitle = "Slot " + itoa(a.selectedSlot) + ": 1=New  0=Back"
		}
	case ModeEnterName:
		vm.MenuTitle = "Enter name (any printable), Enter=OK, Backspace=delete:"
		vm.MenuItems = []string{a.name}
	case ModeConfirmOverwrite:
		vm.MenuTitle = "Overwrite slot " + itoa(a.selectedSlot) + "? (y/n)"
	case ModeConfirmDelete:
		vm.MenuTitle = "Delete save in slot " + itoa(a.selectedSlot) + "? (y/n)"
	case ModeSelectWeapon:
		vm.MenuTitle = "Select weapon (0=unequip, 1-9):"
		vm.MenuItems = weaponsMenu(a.sess.Backpack.List(game.ItemWeapon))
	case ModeSelectFood:
		vm.MenuTitle = "Select food (1-9):"
		vm.MenuItems = itemsMenu(a.sess.Backpack.List(game.ItemFood))
	case ModeSelectElixir:
		vm.MenuTitle = "Select elixir (1-9):"
		vm.MenuItems = itemsMenu(a.sess.Backpack.List(game.ItemElixir))
	case ModeSelectScroll:
		vm.MenuTitle = "Select scroll (1-9):"
		vm.MenuItems = itemsMenu(a.sess.Backpack.List(game.ItemScroll))
	case ModeStats:
		vm.MenuTitle = "Statistics (press any key to close)"
		if rows, err := a.storage.AllRuns(); err == nil {
			vm.StatsLines = statsBoardLines(rows)
		}
	case ModeLeaderboard:
		vm.MenuTitle = "Leaderboard (press any key to close)"
		if rows, err := a.storage.AllRuns(); err == nil {
			vm.BoardLines = boardLines(rows)
		}
	case ModeHelp:
		vm.MenuTitle = "Help (press any key to close)"
		vm.StatsLines = helpLines()
	case ModeGameOver:
		vm.MenuTitle = "Game Over (press any key)"
	case ModeWin:
		vm.MenuTitle = "You Won (press any key)"
	}

	return vm
}

func itemsMenu(items []game.Item) []string {
	out := make([]string, 0, len(items))
	for i := range items {
		out = append(out, items[i].DisplayName())
	}
	return out
}

func weaponsMenu(items []game.Item) []string {
	out := make([]string, 0, len(items))
	for i := range items {
		name := items[i].DisplayName()
		out = append(out, name+" (+"+itoa(items[i].Strength)+")")
	}
	return out
}

func statsLines(s game.Stats) []string {
	return []string{
		"Total treasure: " + itoa(s.TotalTreasure),
		"Deepest level: " + itoa(s.DeepestLevelReached),
		"Defeated enemies: " + itoa(s.DefeatedEnemies),
		"Food consumed: " + itoa(s.FoodConsumed),
		"Elixirs drunk: " + itoa(s.ElixirsDrunk),
		"Scrolls read: " + itoa(s.ScrollsRead),
		"Hits dealt: " + itoa(s.HitsDealt),
		"Hits received: " + itoa(s.HitsReceived),
		"Tiles traveled: " + itoa(s.TilesTraveled),
	}
}

func statsBoardLines(rows []RunResult) []string {
	out := make([]string, 0, len(rows))
	for i := range rows {
		name := rows[i].Name
		if name == "" {
			name = "player"
		}
		s := rows[i].Stats
		out = append(out,
			itoa(i+1)+". "+name+" treasure="+itoa(rows[i].Treasure)+
				" depth="+itoa(rows[i].DeepestLevel)+
				" kills="+itoa(s.DefeatedEnemies)+
				" food="+itoa(s.FoodConsumed)+
				" hits="+itoa(s.HitsDealt)+"/"+itoa(s.HitsReceived)+
				" steps="+itoa(s.TilesTraveled),
		)
	}
	if len(out) == 0 {
		out = append(out, "(no completed runs yet)")
	}
	return out
}

func boardLines(rows []RunResult) []string {
	out := make([]string, 0, len(rows))
	for i := range rows {
		name := rows[i].Name
		if name == "" {
			name = "player"
		}
		out = append(out, itoa(i+1)+". "+name+" treasure="+itoa(rows[i].Treasure)+" depth="+itoa(rows[i].DeepestLevel))
	}
	return out
}

func helpLines() []string {
	return []string{
		"Combat:",
		"- To deal damage: MOVE INTO an enemy tile OR press SPACE to attack adjacent enemy.",
		"- Each of your turns triggers enemy turns.",
		"- Hit chance depends on Dexterity.",
		"- Damage scales with Strength and weapon bonus (+Str).",
		"- XP is earned on kills and levels up your hero (more HP/Str/Dex).",
		"",
		"Weapons:",
		"- Weapons have a +Str bonus shown like Sword (+2).",
		"- Your weapon is an equipped slot referencing a backpack entry.",
		"- When you equip a different weapon, the previously equipped one is DROPPED (removed from backpack).",
		"",
		"Inventory:",
		"- Walk onto items to pick up (if backpack has space).",
		"- h: weapon, j: food, k: elixir, e: scroll (pick 1-9).",
		"- Ctrl+S: save now.",
		"",
		"Enemies:",
		"- z Zombie: high HP.",
		"- v Vampire: first hit always misses; drains Max HP on hit.",
		"- g Ghost: teleports in rooms; sometimes invisible.",
		"- O Ogre: very strong; 2 tiles/turn in rooms; rests after attack then guaranteed hit.",
		"- s Snake-Mage: diagonal movement; may put you to sleep for 1 turn.",
	}
}

func isNameRune(r rune) bool {
	// Accept any printable rune, so users can type names freely (including symbols like !:?)).
	// Reject controls and DEL.
	if r == utf8.RuneError || r == 127 {
		return false
	}
	if unicode.IsControl(r) {
		return false
	}
	return unicode.IsPrint(r)
}

func backspace(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	return string(r[:len(r)-1])
}

func normalizeName(s string) string {
	r := []rune(s)
	for len(r) > 0 && r[0] == ' ' {
		r = r[1:]
	}
	for len(r) > 0 && r[len(r)-1] == ' ' {
		r = r[:len(r)-1]
	}
	return string(r)
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [32]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + (v % 10))
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
