package tui

import (
	"errors"
	"log"

	"rogue/internal/app/usecase"
	"rogue/internal/domain/game"
	"rogue/internal/platform/tcellx"

	"github.com/gdamore/tcell/v2"
)

type Options struct {
	Logger *log.Logger
}

func Run(uc usecase.Controller, opts Options) (err error) {
	if uc == nil {
		return errors.New("nil usecase")
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
	uc.SetWorldSize(worldW, worldH)

	vm, _, err := uc.Handle(usecase.Input{Kind: usecase.InputNone})
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
			if in.Kind == usecase.InputNone {
				continue
			}
			var quit bool
			vm, quit, err = uc.Handle(in)
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

func mapKey(e *tcell.EventKey, mode usecase.Mode) usecase.Input {
	switch e.Key() {
	case tcell.KeyCtrlC, tcell.KeyEscape:
		return usecase.Input{Kind: usecase.InputQuit}
	case tcell.KeyCtrlS:
		return usecase.Input{Kind: usecase.InputSave}
	case tcell.KeyEnter:
		return usecase.Input{Kind: usecase.InputConfirmYes}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return usecase.Input{Kind: usecase.InputBackspace}
	case tcell.KeyRune:
		// During name entry, treat all printable runes as text (do not bind gameplay keys).
		if mode == usecase.ModeEnterName {
			r := e.Rune()
			if r >= 32 && r != 127 {
				return usecase.Input{Kind: usecase.InputText, Rune: r}
			}
			return usecase.Input{Kind: usecase.InputNone}
		}
		switch e.Rune() {
		case 'q', 'Q':
			return usecase.Input{Kind: usecase.InputQuit}
		case 'w', 'W':
			return usecase.Input{Kind: usecase.InputUp}
		case 's', 'S':
			return usecase.Input{Kind: usecase.InputDown}
		case 'a', 'A':
			return usecase.Input{Kind: usecase.InputLeft}
		case 'd', 'D':
			return usecase.Input{Kind: usecase.InputRight}
		case 'h', 'H':
			return usecase.Input{Kind: usecase.InputInvWeapon}
		case 'j', 'J':
			return usecase.Input{Kind: usecase.InputInvFood}
		case 'k', 'K':
			return usecase.Input{Kind: usecase.InputInvElixir}
		case 'e', 'E':
			return usecase.Input{Kind: usecase.InputInvScroll}
		case 't', 'T':
			return usecase.Input{Kind: usecase.InputViewStats}
		case 'l', 'L':
			return usecase.Input{Kind: usecase.InputViewLeaderboard}
		case '?':
			return usecase.Input{Kind: usecase.InputHelp}
		case 'y', 'Y':
			return usecase.Input{Kind: usecase.InputConfirmYes}
		case 'n', 'N':
			return usecase.Input{Kind: usecase.InputConfirmNo}
		case '0':
			return usecase.Input{Kind: usecase.InputDigit0}
		case '1':
			return usecase.Input{Kind: usecase.InputDigit1}
		case '2':
			return usecase.Input{Kind: usecase.InputDigit2}
		case '3':
			return usecase.Input{Kind: usecase.InputDigit3}
		case '4':
			return usecase.Input{Kind: usecase.InputDigit4}
		case '5':
			return usecase.Input{Kind: usecase.InputDigit5}
		case '6':
			return usecase.Input{Kind: usecase.InputDigit6}
		case '7':
			return usecase.Input{Kind: usecase.InputDigit7}
		case '8':
			return usecase.Input{Kind: usecase.InputDigit8}
		case '9':
			return usecase.Input{Kind: usecase.InputDigit9}
		default:
			// Ignore other keys in play modes.
			return usecase.Input{Kind: usecase.InputNone}
		}
	default:
		return usecase.Input{Kind: usecase.InputNone}
	}
}

func render(s tcell.Screen, vm usecase.ViewModel) {
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

func drawMap(s tcell.Screen, ox, oy, mw, mh, sw, sh int, vm usecase.ViewModel) {
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

func drawSidebar(s tcell.Screen, ox, oy, w, h int, vm usecase.ViewModel) {
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

	line("Hero: L"+itoa(vm.Player.Level)+"  XP "+itoa(vm.Player.XP)+"/"+itoa(game.XPToNextLevel(vm.Player.Level)), tcell.StyleDefault)
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
	line("Food    : "+itoa(vm.Backpack.Count(game.ItemFood)), tcell.StyleDefault)
	line("Elixirs : "+itoa(vm.Backpack.Count(game.ItemElixir)), tcell.StyleDefault)
	line("Scrolls : "+itoa(vm.Backpack.Count(game.ItemScroll)), tcell.StyleDefault)
	line("Weapons : "+itoa(vm.Backpack.Count(game.ItemWeapon)), tcell.StyleDefault)
}

func drawLog(s tcell.Screen, ox, oy, w, h int, vm usecase.ViewModel) {
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

func drawModal(s tcell.Screen, sw, sh int, vm usecase.ViewModel) {
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

func tileGlyph(t game.Tile, visible bool, x, y int, vm usecase.ViewModel) (rune, tcell.Style) {
	st := tcell.StyleDefault
	if visible {
		st = st.Foreground(tcell.ColorWhite)
	} else {
		st = st.Foreground(tcell.ColorGray).Dim(true)
	}
	switch t {
	case game.TileWall:
		return wallRune(vm.Tiles, x, y), st.Foreground(tcell.ColorGray)
	case game.TileFloor:
		return '·', st.Foreground(tcell.ColorDarkGray)
	case game.TileCorridor:
		return '∙', st.Foreground(tcell.ColorDarkGray)
	case game.TileExit:
		return '>', st.Foreground(tcell.ColorLightGreen).Bold(true)
	default:
		return ' ', tcell.StyleDefault
	}
}

func wallRune(tiles [][]game.Tile, x, y int) rune {
	// Box-drawing wall based on cardinal neighbors (classic roguelike trick).
	// Treat out-of-bounds as wall for nicer borders.
	isWall := func(xx, yy int) bool {
		if yy < 0 || yy >= len(tiles) || xx < 0 || xx >= len(tiles[yy]) {
			return true
		}
		return tiles[yy][xx] == game.TileWall
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

func drawTopBar(s tcell.Screen, x, y, w int, vm usecase.ViewModel) {
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

func enemyGlyph(e game.Enemy) (rune, tcell.Style) {
	switch e.Type {
	case game.EnemyZombie:
		return 'z', tcell.StyleDefault.Foreground(tcell.ColorGreen)
	case game.EnemyVampire:
		return 'v', tcell.StyleDefault.Foreground(tcell.ColorRed)
	case game.EnemyGhost:
		if e.GhostInvisible {
			return ' ', tcell.StyleDefault
		}
		return 'g', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	case game.EnemyOgre:
		return 'O', tcell.StyleDefault.Foreground(tcell.ColorYellow)
	case game.EnemySnakeMage:
		return 's', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	default:
		return 'm', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	}
}

func itemGlyph(it game.Item) rune {
	switch it.Type {
	case game.ItemFood:
		return ':'
	case game.ItemElixir:
		return '!'
	case game.ItemScroll:
		return '?'
	case game.ItemWeapon:
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
