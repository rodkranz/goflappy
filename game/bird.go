package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"sync"
)

const (
	gravity   = 0.2
	jumpSpeed = 5
)

type bird struct {
	mu sync.RWMutex
	
	time     int
	textures []*sdl.Texture
	
	x, y  int32
	h, w  int32
	speed float64
	dead  bool
}

func NewBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("./res/imgs/birds/player%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load textures image: %v", err)
		}
		textures = append(textures, texture)
	}
	
	return &bird{
		textures: textures,
		x:        10,
		y:        300,
		w:        50,
		h:        43,
	}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) restart() {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	return b.dead
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	rect := &sdl.Rect{X: b.x, Y: 600 - int32(b.y) - b.h/2, W: b.w, H: b.h}
	
	i := b.time / 10 % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy birds: %v", err)
	}
	
	return nil
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.speed = -jumpSpeed
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// too far right
	if p.x > b.x+b.w {
		return
	}
	// too far left
	if p.x+p.w < b.x {
		return
	}
	
	// pipes is too low
	if !p.inverted && p.h < b.y-(b.h/2) {
		return
	}
	
	// pipes is too high
	if p.inverted && 600-p.h > b.y+b.h/2 {
		return
	}
	
	b.dead = true
}
