package game

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type panelPoint struct {
	mu sync.RWMutex

	f *ttf.Font
	c sdl.Color

	points int
}

func newPanelPoint() (_ *panelPoint, err error) {
	pp := new(panelPoint)

	pp.f, err = ttf.OpenFont("./res/fonts/FiraCode.ttf", 20)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	pp.c = sdl.Color{R: 0, G: 0, B: 0, A: 0}
	return pp, nil
}

func (pp *panelPoint) count(ps *pipes) {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pp.points = ps.count
}

func (pp *panelPoint) destroy() {
	pp.f.Close()
}

func (pp *panelPoint) restart() {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	pp.points = 0
}

func (pp *panelPoint) paint(r *sdl.Renderer) error {
	s, err := pp.f.RenderUTF8_Solid(fmt.Sprintf("Points: %d", pp.points), pp.c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	rect := &sdl.Rect{X: 10, Y: 10, W: 100, H: 30}
	if err := r.CopyEx(t, nil, rect, 0, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not copy pipes: %v", err)
	}

	return nil
}
