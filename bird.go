package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"github.com/veandco/go-sdl2/img"
)

type Bird struct {
	time     int
	textures []*sdl.Texture
}

func NewBird(r *sdl.Renderer) (*Bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("./res/imgs/birds/player%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load textures image: %v", err)
		}
		textures = append(textures, texture)
	}
	
	return &Bird{textures: textures}, nil
}

func (b *Bird) paint(r *sdl.Renderer) error {
	b.time++
	
	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}
	
	i := b.time / 10 % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy birds: %v", err)
	}
	
	return nil
}

func (b *Bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}
