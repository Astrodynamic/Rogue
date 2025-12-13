package game

func make2D[T any](w, h int, fill T) [][]T {
	out := make([][]T, h)
	for y := 0; y < h; y++ {
		row := make([]T, w)
		for x := 0; x < w; x++ {
			row[x] = fill
		}
		out[y] = row
	}
	return out
}

func inBounds(p Point, w, h int) bool {
	return p.X >= 0 && p.X < w && p.Y >= 0 && p.Y < h
}
