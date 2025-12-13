package tcellx

import "github.com/gdamore/tcell/v2"

func NewScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	s.SetStyle(tcell.StyleDefault)
	s.Clear()
	return s, nil
}


