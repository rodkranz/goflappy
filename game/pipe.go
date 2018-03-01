package game

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	mu sync.RWMutex

	x int32
	h int32
	w int32

	inverted bool
}

func newPipe() *pipe {
	return &pipe{
		x:        800,
		h:        100 + int32(rand.Intn(300)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE

	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipes: %v", err)
	}

	return nil
}

func (p *pipe) setSpeed(s int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= s
}

func (p *pipe) touch(b *bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	b.touch(p)
}
