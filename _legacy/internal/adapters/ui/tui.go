package ui

// UI imports domain types (Tile, Enemy, Point, etc.) used in ViewModel.
// In strict Clean Architecture this would be a violation (adapter knows domain).
// In Go this is a pragmatic tradeoff: duplicating 15+ value types just for
// isolation adds complexity without real benefit. The key rule is preserved:
// UI never CALLS domain logic, only reads data through ViewModel.

import (
	"errors"
	"log"

	"rogue/internal/service"
	"rogue/internal/domain"
	"rogue/internal/platform/tcellx"

	"github.com/gdamore/tcell/v2"
)

type Options struct {
	Logger *log.Logger
}

func Run(ctrl service.Controller, opts Options) (err error) {
	if ctrl == nil {
		return errors.New("nil controller")
	}
	if opts.Logger == nil {
		opts.Logger = log.Default()
	}

	s, err := tcellx.NewScreen()
	if err != nil {
		return err
	}
	defer func() {
		s.Fini()
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	// Presentation layer owns terminal concerns: compute initial world size from screen.
	sw, sh := s.Size()
	// Account for UI chrome: sidebar + borders + top bar.
	worldW := sw - 35
	worldH := sh - 6
	ctrl.SetWorldSize(worldW, worldH)

	vm, _, err := ctrl.Handle(service.Input{Kind: service.InputNone})
	if err != nil {
		return err
	}
	render(s, vm)

	for {
		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			render(s, vm)
		case *tcell.EventKey:
			in := mapKey(e, vm.Mode)
			if in.Kind == service.InputNone {
				continue
			}
			var quit bool
			vm, quit, err = ctrl.Handle(in)
			if err != nil {
				return err
			}
			render(s, vm)
			if quit {
				return nil
			}
		}
	}
}

func mapKey(e *tcell.EventKey, mode service.Mode) service.Input {
	switch e.Key() {
	case tcell.KeyCtrlC, tcell.KeyEscape:
		return service.Input{Kind: service.InputQuit}
	case tcell.KeyCtrlS:
		return service.Input{Kind: service.InputSave}
	case tcell.KeyEnter:
		return service.Input{Kind: service.InputConfirmYes}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return service.Input{Kind: service.InputBackspace}
	case tcell.KeyRune:
		// During name entry, treat all printable runes as text (do not bind gameplay keys).
		if mode == service.ModeEnterName {
			r := e.Rune()
			if r >= 32 && r != 127 {
				return service.Input{Kind: service.InputText, Rune: r}
			}
			return service.Input{Kind: service.InputNone}
		}
		switch e.Rune() {
		case ' ':
			return service.Input{Kind: service.InputAttack}
		case 'q', 'Q':
			return service.Input{Kind: service.InputQuit}
		case 'w', 'W':
			return service.Input{Kind: service.InputUp}
		case 's', 'S':
			return service.Input{Kind: service.InputDown}
		case 'a', 'A':
			return service.Input{Kind: service.InputLeft}
		case 'd', 'D':
			return service.Input{Kind: service.InputRight}
		case 'h', 'H':
			return service.Input{Kind: service.InputInvWeapon}
		case 'j', 'J':
			return service.Input{Kind: service.InputInvFood}
		case 'k', 'K':
			return service.Input{Kind: service.InputInvElixir}
		case 'e', 'E':
			return service.Input{Kind: service.InputInvScroll}
		case 't', 'T':
			return service.Input{Kind: service.InputViewStats}
		case 'l', 'L':
			return service.Input{Kind: service.InputViewLeaderboard}
		case '?':
			return service.Input{Kind: service.InputHelp}
		case 'y', 'Y':
			return service.Input{Kind: service.InputConfirmYes}
		case 'n', 'N':
			return service.Input{Kind: service.InputConfirmNo}
		case '0':
			return service.Input{Kind: service.InputDigit0}
		case '1':
			return service.Input{Kind: service.InputDigit1}
		case '2':
			return service.Input{Kind: service.InputDigit2}
		case '3':
			return service.Input{Kind: service.InputDigit3}
		case '4':
			return service.Input{Kind: service.InputDigit4}
		case '5':
			return service.Input{Kind: service.InputDigit5}
		case '6':
			return service.Input{Kind: service.InputDigit6}
		case '7':
			return service.Input{Kind: service.InputDigit7}
		case '8':
			return service.Input{Kind: service.InputDigit8}
		case '9':
			return service.Input{Kind: service.InputDigit9}
		default:
			// Ignore other keys in play modes.
			return service.Input{Kind: service.InputNone}
		}
	default:
		return service.Input{Kind: service.InputNone}
	}
}

func render(s tcell.Screen, vm service.ViewModel) {
	s.Clear()
	w, h := s.Size()
	if w <= 0 || h <= 0 {
		s.Show()
		return
	}

	topH := 2
	if h < topH+6 || w < 50 {
		drawText(s, 0, 0, tcell.StyleDefault.Bold(true), truncate(vm.Title, w))
		drawText(s, 0, 1, tcell.StyleDefault, truncate(vm.Help, w))
		s.Show()
		return
	}

	drawTopBar(s, 0, 0, w, vm)

	sideW := 30
	if w < vm.MapW+sideW+6 {
		sideW = maxInt(18, w-vm.MapW-6)
	}
	mapW := minInt(vm.MapW, w-sideW-3)
	mapH := minInt(vm.MapH, h-topH-2)
	mapX, mapY := 0, topH
	sideX := mapX + mapW + 1
	sideY := topH
	statusH := minInt(12, maxInt(8, h-topH-6))
	logY := sideY + statusH

	drawBox(s, mapX, mapY, mapW+2, mapH+2, "Dungeon", tcell.ColorGray)
	drawMap(s, mapX+1, mapY+1, mapW, mapH, w, h, vm)

	drawBox(s, sideX, sideY, sideW, statusH, "Status", tcell.ColorGray)
	drawSidebar(s, sideX+1, sideY+1, sideW-2, statusH-2, vm)

	logH := h - logY - 1
	if logH >= 4 {
		drawBox(s, sideX, logY, sideW, logH, "Log", tcell.ColorGray)
		drawLog(s, sideX+1, logY+1, sideW-2, logH-2, vm)
	}

	if vm.MenuTitle != "" {
		drawModal(s, w, h, vm)
	}

	s.Show()
}

func drawText(s tcell.Screen, x, y int, st tcell.Style, txt string) {
	for i, r := range txt {
		s.SetContent(x+i, y, r, nil, st)
	}
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

func drawMap(s tcell.Screen, ox, oy, mw, mh, sw, sh int, vm service.ViewModel) {
	roomTile := make([][]bool, vm.MapH)
	for y := 0; y < vm.MapH; y++ {
		roomTile[y] = make([]bool, vm.MapW)
	}
	for _, r := range vm.Rooms {
		for y := r.Bounds.Y; y < r.Bounds.Y+r.Bounds.H; y++ {
			if y < 0 || y >= vm.MapH {
				continue
			}
			for x := r.Bounds.X; x < r.Bounds.X+r.Bounds.W; x++ {
				if x < 0 || x >= vm.MapW {
					continue
				}
				roomTile[y][x] = true
			}
		}
	}

	enemyAt := make(map[[2]int]rune, len(vm.Enemies))
	enemyStyle := make(map[[2]int]tcell.Style, len(vm.Enemies))
	for _, e := range vm.Enemies {
		k := [2]int{e.Pos.X, e.Pos.Y}
		ch, st := enemyGlyph(e)
		enemyAt[k] = ch
		enemyStyle[k] = st
	}
	itemAt := make(map[[2]int]rune, 32)
	for _, r := range vm.Rooms {
		for p, it := range r.FloorLoot {
			itemAt[[2]int{p.X, p.Y}] = itemGlyph(it)
		}
	}

	// Camera/viewport centered on the player.
	vx0 := clampInt(vm.PlayerPos.X-mw/2, 0, maxInt(0, vm.MapW-mw))
	vy0 := clampInt(vm.PlayerPos.Y-mh/2, 0, maxInt(0, vm.MapH-mh))

	stDim := tcell.StyleDefault.Foreground(tcell.ColorGray).Dim(true)
	stVisible := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	for vy := 0; vy < mh; vy++ {
		y := vy0 + vy
		if y < 0 || y >= vm.MapH || oy+vy >= sh {
			continue
		}
		for vx := 0; vx < mw; vx++ {
			x := vx0 + vx
			if x < 0 || x >= vm.MapW || ox+vx >= sw {
				continue
			}
			if !vm.Explored[y][x] {
				continue
			}
			visible := vm.Visible[y][x]

			ch, st := tileGlyph(vm.Tiles[y][x], visible, x, y, vm)
			if !visible {
				// Spec: explored rooms/corridors where player is not present are displayed as walls only.
				ch, st = wallRune(vm.Tiles, x, y), stDim
				_ = roomTile // kept for potential later fog nuance
			}

			if visible {
				if x == vm.PlayerPos.X && y == vm.PlayerPos.Y {
					ch = '@'
					st = stVisible.Bold(true)
				} else if eg, ok := enemyAt[[2]int{x, y}]; ok {
					ch = eg
					st = enemyStyle[[2]int{x, y}]
				} else if ig, ok := itemAt[[2]int{x, y}]; ok {
					ch = ig
					st = tcell.StyleDefault.Foreground(tcell.ColorLightCyan).Bold(true)
				}
			}
			s.SetContent(ox+vx, oy+vy, ch, nil, st)
		}
	}
}

func drawSidebar(s tcell.Screen, ox, oy, w, h int, vm service.ViewModel) {
	y := oy
	line := func(txt string, st tcell.Style) {
		if h <= 0 {
			return
		}
		drawText(s, ox, y, st, truncate(txt, w))
		y++
		h--
	}

	line("Dungeon: "+itoa(vm.Level), tcell.StyleDefault.Bold(true))
	if h <= 0 {
		return
	}
	drawHPBar(s, ox, y, w, vm.Player.Health, vm.Player.MaxHealth)
	y++
	h--

	line("Hero: L"+itoa(vm.Player.Level)+"  XP "+itoa(vm.Player.XP)+"/"+itoa(domain.XPToNextLevel(vm.Player.Level)), tcell.StyleDefault)
	line("Dex: "+itoa(vm.Player.Dexterity)+"  Str: "+itoa(vm.Player.Strength), tcell.StyleDefault)
	line("Treasure: "+itoa(vm.Backpack.Treasure), tcell.StyleDefault)

	weapon := "none"
	if vm.Weapon != nil {
		weapon = vm.Weapon.DisplayName() + " (+" + itoa(vm.Weapon.Strength) + ")"
	}
	line("Weapon: "+weapon, tcell.StyleDefault)

	if h > 0 {
		y++
		h--
	}
	line("Food    : "+itoa(vm.Backpack.Count(domain.ItemFood)), tcell.StyleDefault)
	line("Elixirs : "+itoa(vm.Backpack.Count(domain.ItemElixir)), tcell.StyleDefault)
	line("Scrolls : "+itoa(vm.Backpack.Count(domain.ItemScroll)), tcell.StyleDefault)
	line("Weapons : "+itoa(vm.Backpack.Count(domain.ItemWeapon)), tcell.StyleDefault)
}

func drawLog(s tcell.Screen, ox, oy, w, h int, vm service.ViewModel) {
	if h <= 0 || w <= 0 {
		return
	}
	msgs := vm.Messages
	if len(msgs) > h {
		msgs = msgs[len(msgs)-h:]
	}
	for i := 0; i < len(msgs) && i < h; i++ {
		drawText(s, ox, oy+i, tcell.StyleDefault.Foreground(tcell.ColorWhite), truncate(msgs[i], w))
	}
}

func drawModal(s tcell.Screen, sw, sh int, vm service.ViewModel) {
	lines := []string{vm.MenuTitle}
	if len(vm.MenuItems) > 0 {
		for i := range vm.MenuItems {
			lines = append(lines, itoa(i+1)+") "+vm.MenuItems[i])
		}
	}
	if len(vm.StatsLines) > 0 {
		lines = append(lines, vm.StatsLines...)
	}
	if len(vm.BoardLines) > 0 {
		lines = append(lines, vm.BoardLines...)
	}

	maxW := 0
	for _, ln := range lines {
		if l := len([]rune(ln)); l > maxW {
			maxW = l
		}
	}
	boxW := maxW + 4
	if boxW > sw-2 {
		boxW = sw - 2
	}
	boxH := len(lines) + 2
	if boxH > sh-2 {
		boxH = sh - 2
	}
	x0 := (sw - boxW) / 2
	y0 := (sh - boxH) / 2

	drawBox(s, x0, y0, boxW, boxH, "", tcell.ColorLightBlue)

	for i := 0; i < len(lines) && i < boxH-2; i++ {
		st := tcell.StyleDefault
		if i == 0 {
			st = st.Bold(true)
		}
		drawText(s, x0+2, y0+1+i, st, truncate(lines[i], boxW-4))
	}
}

func tileGlyph(t domain.Tile, visible bool, x, y int, vm service.ViewModel) (rune, tcell.Style) {
	st := tcell.StyleDefault
	if visible {
		st = st.Foreground(tcell.ColorWhite)
	} else {
		st = st.Foreground(tcell.ColorGray).Dim(true)
	}
	switch t {
	case domain.TileWall:
		return wallRune(vm.Tiles, x, y), st.Foreground(tcell.ColorGray)
	case domain.TileFloor:
		return '·', st.Foreground(tcell.ColorDarkGray)
	case domain.TileCorridor:
		return '∙', st.Foreground(tcell.ColorDarkGray)
	case domain.TileExit:
		return '>', st.Foreground(tcell.ColorLightGreen).Bold(true)
	default:
		return ' ', tcell.StyleDefault
	}
}

func wallRune(tiles [][]domain.Tile, x, y int) rune {
	// Box-drawing wall based on cardinal neighbors (classic roguelike trick).
	// Treat out-of-bounds as wall for nicer borders.
	isWall := func(xx, yy int) bool {
		if yy < 0 || yy >= len(tiles) || xx < 0 || xx >= len(tiles[yy]) {
			return true
		}
		return tiles[yy][xx] == domain.TileWall
	}
	up := isWall(x, y-1)
	down := isWall(x, y+1)
	left := isWall(x-1, y)
	right := isWall(x+1, y)

	switch {
	case up && down && left && right:
		return '┼'
	case up && down && left:
		return '┤'
	case up && down && right:
		return '├'
	case left && right && up:
		return '┴'
	case left && right && down:
		return '┬'
	case up && down:
		return '│'
	case left && right:
		return '─'
	case down && right:
		return '┌'
	case down && left:
		return '┐'
	case up && right:
		return '└'
	case up && left:
		return '┘'
	case up || down:
		return '│'
	case left || right:
		return '─'
	default:
		return '█'
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func drawTopBar(s tcell.Screen, x, y, w int, vm service.ViewModel) {
	st := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true)
	fill(s, x, y, w, 1, ' ', st)
	title := vm.Title + "  |  Level " + itoa(vm.Level) + "  |  Treasure " + itoa(vm.Backpack.Treasure)
	drawText(s, x+1, y, st, truncate(title, w-2))

	st2 := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray)
	fill(s, x, y+1, w, 1, ' ', st2)
	drawText(s, x+1, y+1, st2, truncate(vm.Help, w-2))
}

func drawBox(s tcell.Screen, x, y, w, h int, title string, borderColor tcell.Color) {
	if w < 2 || h < 2 {
		return
	}
	st := tcell.StyleDefault.Foreground(borderColor)
	s.SetContent(x, y, tcell.RuneULCorner, nil, st)
	s.SetContent(x+w-1, y, tcell.RuneURCorner, nil, st)
	s.SetContent(x, y+h-1, tcell.RuneLLCorner, nil, st)
	s.SetContent(x+w-1, y+h-1, tcell.RuneLRCorner, nil, st)
	for i := 1; i < w-1; i++ {
		s.SetContent(x+i, y, tcell.RuneHLine, nil, st)
		s.SetContent(x+i, y+h-1, tcell.RuneHLine, nil, st)
	}
	for j := 1; j < h-1; j++ {
		s.SetContent(x, y+j, tcell.RuneVLine, nil, st)
		s.SetContent(x+w-1, y+j, tcell.RuneVLine, nil, st)
	}
	if title != "" && w > 4 {
		drawText(s, x+2, y, st.Bold(true), truncate(title, w-4))
	}
}

func fill(s tcell.Screen, x, y, w, h int, r rune, st tcell.Style) {
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			s.SetContent(x+xx, y+yy, r, nil, st)
		}
	}
}

func drawHPBar(s tcell.Screen, x, y, w, hp, maxHP int) {
	if w <= 0 || maxHP <= 0 {
		return
	}
	label := "HP " + itoa(hp) + "/" + itoa(maxHP) + " "
	drawText(s, x, y, tcell.StyleDefault, truncate(label, w))
	barW := w - len([]rune(label))
	if barW <= 0 {
		return
	}
	if hp < 0 {
		hp = 0
	}
	if hp > maxHP {
		hp = maxHP
	}
	filled := (hp * barW) / maxHP
	col := tcell.ColorGreen
	if hp*100/maxHP < 30 {
		col = tcell.ColorRed
	} else if hp*100/maxHP < 60 {
		col = tcell.ColorYellow
	}
	for i := 0; i < barW; i++ {
		st := tcell.StyleDefault
		ch := '░'
		if i < filled {
			ch = '█'
			st = st.Foreground(col)
		} else {
			st = st.Foreground(tcell.ColorGray).Dim(true)
		}
		s.SetContent(x+len([]rune(label))+i, y, ch, nil, st)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func enemyGlyph(e domain.Enemy) (rune, tcell.Style) {
	switch e.Type {
	case domain.EnemyZombie:
		return 'z', tcell.StyleDefault.Foreground(tcell.ColorGreen)
	case domain.EnemyVampire:
		return 'v', tcell.StyleDefault.Foreground(tcell.ColorRed)
	case domain.EnemyGhost:
		if e.GhostInvisible {
			return ' ', tcell.StyleDefault
		}
		return 'g', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	case domain.EnemyOgre:
		return 'O', tcell.StyleDefault.Foreground(tcell.ColorYellow)
	case domain.EnemySnakeMage:
		return 's', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	default:
		return 'm', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	}
}

func itemGlyph(it domain.Item) rune {
	switch it.Type {
	case domain.ItemFood:
		return ':'
	case domain.ItemElixir:
		return '!'
	case domain.ItemScroll:
		return '?'
	case domain.ItemWeapon:
		return ')'
	default:
		return '*'
	}
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
