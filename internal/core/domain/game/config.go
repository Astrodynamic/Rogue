package game

type Config struct {
	Width  int
	Height int
}

func (c Config) Normalize() Config {
	// Minimums to keep 3x3 room grid workable.
	if c.Width < 60 {
		c.Width = 60
	}
	if c.Height < 18 {
		c.Height = 18
	}
	// Upper bounds keep UI sane and generation fast.
	if c.Width > 140 {
		c.Width = 140
	}
	if c.Height > 45 {
		c.Height = 45
	}
	return c
}
