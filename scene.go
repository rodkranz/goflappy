package main

import (
	"fmt"
	
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"context"
	"time"
)

type scene struct {
	time int
	
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "./res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	
	var birds []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("./res/imgs/birds/player%d.png", i)
		bird, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load birds image: %v", err)
		}
		birds = append(birds, bird)
	}
	
	return &scene{bg: bg, birds: birds}, nil
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	
	return errc
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	
	r.Clear()
	
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	
	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}
	
	i := s.time / 10 % len(s.birds)
	if err := r.Copy(s.birds[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy birds: %v", err)
	}
	
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
