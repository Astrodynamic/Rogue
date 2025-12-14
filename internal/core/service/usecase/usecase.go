package usecase

import (
	"errors"
	"math/rand"
	"time"
	"unicode"
	"unicode/utf8"

	"rogue/internal/core/domain/game"
	"rogue/internal/core/ports"
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

type Usecase struct {
	svc App

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
func (u *Usecase) SetWorldSize(width, height int) {
	u.cfg = game.Config{Width: width, Height: height}.Normalize()
	// If we're still in slot selection (pre-run), keep the seed but regenerate the initial level for a nicer first impression.
	if u.mode == ModeSelectSlot && u.sess != nil {
		u.sess = game.NewSession(u.sess.Seed, u.cfg)
	}
}

func New(svc App, seed int64, cfg game.Config) *Usecase {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))
	return &Usecase{
		svc:  svc,
		rng:  rng,
		sess: game.NewSession(seed, cfg),
		mode: ModeSelectSlot,
		cfg:  cfg.Normalize(),
	}
}

func (u *Usecase) Handle(in Input) (ViewModel, bool, error) {
	if u.sess == nil {
		return ViewModel{}, true, errors.New("nil session")
	}

	switch u.mode {
	case ModeSelectSlot:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		if slot, ok := digitIndex(in); ok && slot >= 1 && slot <= 3 {
			u.selectedSlot = slot
			_, has, _ := u.svc.LoadSession(slot)
			u.selectedHasSave = has
			u.mode = ModeSlotAction
			return u.view(), false, nil
		}
		return u.view(), false, nil
	case ModeSlotAction:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		// 1=Continue/New, 2=New (overwrite), 3=Delete, 0=Back
		if d, ok := digitIndex(in); ok {
			switch d {
			case 0:
				u.mode = ModeSelectSlot
				return u.view(), false, nil
			case 1:
				if u.selectedHasSave {
					u.activeSlot = u.selectedSlot
					if saved, ok2, err := u.svc.LoadSession(u.activeSlot); err == nil && ok2 {
						u.sess = game.FromSnapshot(saved.Snapshot)
						u.rng = rand.New(rand.NewSource(u.sess.Seed))
						u.name = saved.Name
						if normalizeName(u.name) == "" {
							u.mode = ModeEnterName
						} else {
							u.mode = ModePlay
						}
						u.sess.Messages = append(u.sess.Messages, "Loaded slot "+itoa(u.activeSlot)+".")
						return u.view(), false, nil
					}
					u.mode = ModeEnterName
					return u.view(), false, nil
				}
				// Empty slot -> start new
				u.activeSlot = u.selectedSlot
				u.mode = ModeEnterName
				return u.view(), false, nil
			case 2:
				if u.selectedHasSave {
					u.mode = ModeConfirmOverwrite
					return u.view(), false, nil
				}
				u.activeSlot = u.selectedSlot
				u.mode = ModeEnterName
				return u.view(), false, nil
			case 3:
				if u.selectedHasSave {
					u.mode = ModeConfirmDelete
					return u.view(), false, nil
				}
				return u.view(), false, nil
			}
		}
		return u.view(), false, nil
	case ModeConfirmOverwrite:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		if in.Kind == InputConfirmNo {
			u.mode = ModeSlotAction
			return u.view(), false, nil
		}
		if in.Kind == InputConfirmYes {
			u.activeSlot = u.selectedSlot
			_ = u.svc.ClearSession(u.activeSlot)
			u.selectedHasSave = false
			u.mode = ModeEnterName
			return u.view(), false, nil
		}
		return u.view(), false, nil
	case ModeConfirmDelete:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		if in.Kind == InputConfirmNo {
			u.mode = ModeSlotAction
			return u.view(), false, nil
		}
		if in.Kind == InputConfirmYes {
			_ = u.svc.ClearSession(u.selectedSlot)
			u.selectedHasSave = false
			u.mode = ModeSelectSlot
			return u.view(), false, nil
		}
		return u.view(), false, nil
	case ModeStats:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		u.mode = ModePlay
		return u.view(), false, nil
	case ModeLeaderboard:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		u.mode = ModePlay
		return u.view(), false, nil
	case ModeHelp:
		if in.Kind == InputQuit {
			return u.view(), true, nil
		}
		u.mode = ModePlay
		return u.view(), false, nil
	case ModeEnterName:
		switch in.Kind {
		case InputQuit:
			return u.view(), true, nil
		case InputConfirmYes:
			u.name = normalizeName(u.name)
			if u.name == "" {
				u.name = "player"
			}
			if u.activeSlot == 0 {
				u.activeSlot = 1
			}
			seed := time.Now().UnixNano()
			u.rng = rand.New(rand.NewSource(seed))
			u.sess = game.NewSession(seed, u.cfg)
			u.mode = ModePlay
			return u.view(), false, nil
		case InputBackspace:
			u.name = backspace(u.name)
			return u.view(), false, nil
		case InputText:
			if len([]rune(u.name)) < 16 && isNameRune(in.Rune) {
				u.name += string(in.Rune)
			}
			return u.view(), false, nil
		default:
			return u.view(), false, nil
		}
	case ModeSelectWeapon, ModeSelectFood, ModeSelectElixir, ModeSelectScroll:
		idx, ok := digitIndex(in)
		if !ok {
			if in.Kind == InputQuit || in.Kind == InputConfirmNo {
				u.mode = ModePlay
			}
			return u.view(), in.Kind == InputQuit, nil
		}
		u.applySelection(idx)
		u.mode = ModePlay
		return u.view(), false, nil
	case ModeGameOver:
		// Any key restarts.
		u.startNewRun()
		return u.view(), false, nil
	case ModeWin:
		// Any key exits.
		return u.view(), true, nil
	}

	// Play mode.
	switch in.Kind {
	case InputQuit:
		u.saveNow(false)
		return u.view(), true, nil
	case InputSave:
		u.saveNow(true)
	case InputHelp:
		u.mode = ModeHelp
	case InputUp:
		u.step(game.Action{Kind: game.ActionMove, Dx: 0, Dy: -1})
	case InputDown:
		u.step(game.Action{Kind: game.ActionMove, Dx: 0, Dy: 1})
	case InputLeft:
		u.step(game.Action{Kind: game.ActionMove, Dx: -1, Dy: 0})
	case InputRight:
		u.step(game.Action{Kind: game.ActionMove, Dx: 1, Dy: 0})
	case InputAttack:
		u.step(game.Action{Kind: game.ActionAttack})
	case InputInvWeapon:
		u.mode = ModeSelectWeapon
		u.pending = game.ActionEquipWeapon
	case InputInvFood:
		u.mode = ModeSelectFood
		u.pending = game.ActionUseFood
	case InputInvElixir:
		u.mode = ModeSelectElixir
		u.pending = game.ActionUseElixir
	case InputInvScroll:
		u.mode = ModeSelectScroll
		u.pending = game.ActionUseScroll
	case InputViewStats:
		u.mode = ModeStats
	case InputViewLeaderboard:
		u.mode = ModeLeaderboard
	}

	return u.view(), false, nil
}

func (u *Usecase) step(a game.Action) {
	prevLevel := u.sess.LevelDepth
	advanced, dead, won, _ := u.sess.Step(u.rng, a)
	if advanced && u.sess.LevelDepth != prevLevel {
		// Autosave only on level transitions (reaching the exit).
		u.saveNow(false)
	}
	if won {
		u.finishRun(true)
		u.mode = ModeWin
	}
	if dead {
		u.finishRun(false)
		u.mode = ModeGameOver
	}
}

func (u *Usecase) saveNow(showMessage bool) {
	if u.svc == nil || u.sess == nil {
		return
	}
	name := normalizeName(u.name)
	if name == "" {
		// Don't create anonymous saves; name entry flow happens first anyway.
		return
	}
	if u.activeSlot == 0 {
		u.activeSlot = 1
	}
	_ = u.svc.SaveSession(u.activeSlot, ports.SavedSession{Snapshot: u.sess.Snapshot(), Name: name})
	if showMessage {
		u.sess.Messages = append(u.sess.Messages, "Saved (slot "+itoa(u.activeSlot)+").")
	}
}

func (u *Usecase) finishRun(won bool) {
	name := normalizeName(u.name)
	if name == "" {
		name = "player"
	}
	_ = u.svc.RecordRun(ports.RunResult{
		Name:         name,
		DeepestLevel: u.sess.Stats.DeepestLevelReached,
		Treasure:     u.sess.Stats.TotalTreasure,
		Stats:        u.sess.Stats,
	})
	if won {
		u.sess.Messages = append(u.sess.Messages, "You won! Press any key to exit.")
	} else {
		u.sess.Messages = append(u.sess.Messages, "Game over. Press any key to restart.")
	}
}

func (u *Usecase) startNewRun() {
	u.mode = ModeSelectSlot
	u.activeSlot = 0
}

func (u *Usecase) applySelection(idx int) {
	switch u.pending {
	case game.ActionEquipWeapon:
		if idx == 0 {
			u.step(game.Action{Kind: game.ActionUnequipWeapon})
			return
		}
		u.step(game.Action{Kind: game.ActionEquipWeapon, Idx: idx - 1})
	case game.ActionUseFood:
		u.step(game.Action{Kind: game.ActionUseFood, Idx: idx - 1})
	case game.ActionUseElixir:
		u.step(game.Action{Kind: game.ActionUseElixir, Idx: idx - 1})
	case game.ActionUseScroll:
		u.step(game.Action{Kind: game.ActionUseScroll, Idx: idx - 1})
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

func (u *Usecase) view() ViewModel {
	vis := u.sess.ComputeVisible()
	var wp *game.Item
	if wid := u.sess.Player.Equipment.Get(game.SlotWeapon); wid != 0 {
		if it, ok := u.sess.Backpack.FindByID(wid); ok && it.Type == game.ItemWeapon {
			cp := it
			wp = &cp
		}
	}
	mw := len(u.sess.Level.Tiles[0])
	mh := len(u.sess.Level.Tiles)
	vm := ViewModel{
		Title:     "Rogue",
		Help:      "Move WASD  Attack Space  Inv h/j/k/e  Save Ctrl+S  (auto on exit+quit)  Stats t  Board l  Help ?  Quit q",
		Mode:      u.mode,
		Level:     u.sess.LevelDepth,
		MapW:      mw,
		MapH:      mh,
		Tiles:     u.sess.Level.Tiles,
		Explored:  u.sess.Explored,
		Visible:   vis,
		PlayerPos: u.sess.PlayerPos,
		Player:    u.sess.Player,
		Weapon:    wp,
		Backpack:  u.sess.Backpack,
		Enemies:   append([]game.Enemy(nil), u.sess.Enemies...),
		Rooms:     u.sess.Level.Rooms,
		Messages:  append([]string(nil), u.sess.Messages...),
	}

	switch u.mode {
	case ModeSelectSlot:
		vm.MenuTitle = "Select save slot (1-3). Empty slot starts a new game."
		if slots, err := u.svc.ListSaveSlots(); err == nil {
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
		if u.selectedHasSave {
			vm.MenuTitle = "Slot " + itoa(u.selectedSlot) + ": 1=Continue  2=New(overwrite)  3=Delete  0=Back"
		} else {
			vm.MenuTitle = "Slot " + itoa(u.selectedSlot) + ": 1=New  0=Back"
		}
	case ModeEnterName:
		vm.MenuTitle = "Enter name (any printable), Enter=OK, Backspace=delete:"
		vm.MenuItems = []string{u.name}
	case ModeConfirmOverwrite:
		vm.MenuTitle = "Overwrite slot " + itoa(u.selectedSlot) + "? (y/n)"
	case ModeConfirmDelete:
		vm.MenuTitle = "Delete save in slot " + itoa(u.selectedSlot) + "? (y/n)"
	case ModeSelectWeapon:
		vm.MenuTitle = "Select weapon (0=unequip, 1-9):"
		vm.MenuItems = weaponsMenu(u.sess.Backpack.List(game.ItemWeapon))
	case ModeSelectFood:
		vm.MenuTitle = "Select food (1-9):"
		vm.MenuItems = itemsMenu(u.sess.Backpack.List(game.ItemFood))
	case ModeSelectElixir:
		vm.MenuTitle = "Select elixir (1-9):"
		vm.MenuItems = itemsMenu(u.sess.Backpack.List(game.ItemElixir))
	case ModeSelectScroll:
		vm.MenuTitle = "Select scroll (1-9):"
		vm.MenuItems = itemsMenu(u.sess.Backpack.List(game.ItemScroll))
	case ModeStats:
		vm.MenuTitle = "Statistics (press any key to close)"
		if rows, err := u.svc.AllRuns(); err == nil {
			vm.StatsLines = statsBoardLines(rows)
		}
	case ModeLeaderboard:
		vm.MenuTitle = "Leaderboard (press any key to close)"
		if rows, err := u.svc.AllRuns(); err == nil {
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

func statsBoardLines(rows []ports.RunResult) []string {
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

func boardLines(rows []ports.RunResult) []string {
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
