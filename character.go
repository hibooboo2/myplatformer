package main

import (
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	_ Entity = &Player{}
)

type Player struct {
	sync.RWMutex
	texture      *sdl.Texture
	X, Y         int
	H, W         int32
	size         int32
	dir          int32
	frame        int32
	frameCounter int32
	speed        int32
}

func NewPlayer(r *sdl.Renderer) (*Player, error) {
	p := new(Player)

	t, err := img.LoadTexture(r, "assets/char.png")
	if err != nil {
		return nil, err
	}
	p.texture = t
	_, _, w, h, err := t.Query()
	if err != nil {
		return nil, err
	}
	p.H = h / 4
	p.W = w / 4
	return p, nil
}

func (p *Player) Paint(r *sdl.Renderer) error {
	v := r.GetViewport()
	p.RLock()
	defer p.RUnlock()
	return r.Copy(p.texture, p.getFrame(), &sdl.Rect{X: v.W/2 - p.W/2, Y: v.H/2 - p.H/2, W: p.size, H: p.size})
}

func (p *Player) Update() {
	p.Lock()
	p.frameCounter++
	if p.frameCounter%(FramesPerSecond/10) == 0 {
		p.frame = (p.frame + 1) % 4
		if p.frame%4 == 0 {
			p.dir = (p.dir + 1) % 4
		}
	}
	if p.frameCounter > 10000 {
		p.frameCounter = 0
	}
	p.Unlock()
}

func (p *Player) getFrame() *sdl.Rect {
	p.RLock()
	r := &sdl.Rect{W: p.W, H: p.H, X: p.frame * p.W, Y: p.dir * p.H}
	p.RUnlock()
	return r
}

func (p *Player) Destroy() {

}

func (p *Player) Reset() {

}

func (p *Player) Resize(delta int32) {
	p.Lock()
	defer p.Unlock()
	p.size += delta
	if p.size > 256 {
		p.size = 256
	}
	if p.size < 10 {
		p.size = 10
	}
}
