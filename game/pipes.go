package game

import (
	"sync"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"time"
)

type pipes struct {
	mu sync.RWMutex
	
	texture *sdl.Texture
	speed   int32
	
	pipes []*pipe
	
	count int
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(r, "./res/imgs/bpipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipes image: %v", err)
	}
	ps := &pipes{
		texture: texture,
		speed:   2,
	}
	
	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			//ps.count++
			ps.mu.Unlock()
			time.Sleep(1500 * time.Millisecond)
		}
	}()
	
	return ps, nil
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	
	for _, p := range ps.pipes {
		err := p.paint(r, ps.texture)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	ps.pipes = nil
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	
	for _, p := range ps.pipes {
		p.touch(b)
	}
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	var ram []*pipe
	for _, p := range ps.pipes {
		p.setSpeed(ps.speed)
		if p.x+p.w > 0 {
			ram = append(ram, p)
		} else {
			ps.count++
		}
	}
	
	ps.pipes = ram
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	ps.texture.Destroy()
}
