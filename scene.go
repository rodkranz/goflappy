package main

import (
	"fmt"
	"time"
	
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/ttf"
)

type Object interface {
	paint(*sdl.Renderer) error
	restart()
	update()
	destroy()
}

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
	
	panel *panelPoint
	
	points int
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "./res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	
	b, err := NewBird(r)
	if err != nil {
		return nil, err
	}
	
	ps, err := newPipes(r)
	if err != nil {
		return nil, err
	}
	
	pp, err := newPanelPoint(r)
	if err != nil {
		return nil, err
	}
	
	return &scene{
		bg:    bg,
		bird:  b,
		pipes: ps,
		panel: pp,
	}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		
		for {
			select {
			case e := <-events:
				if s.handleEvent(e) {
					return
				}
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(time.Second * 2)
					s.restart()
				}
				
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	}
	
	return false
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
	s.panel.count(s.pipes)
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	if err := s.bird.paint(r); err != nil {
		return err
	}
	if err := s.pipes.paint(r); err != nil {
		return err
	}
	
	if err := s.panel.paint(r); err != nil {
		return err
	}
	
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
	s.panel.destroy()
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()
	
	f, err := ttf.OpenFont("./res/fonts/FiraCode.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()
	
	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8_Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()
	
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()
	
	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}
	
	r.Present()
	return nil
}
