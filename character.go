package main

import (
	"fmt"
	"sync/atomic"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	_ Entity = &Player{}
)

type Player struct {
	texture      *sdl.Texture
	X, Y         *int
	H, W         *int32
	size         *int32
	dir          *int32
	dirx         *int32
	diry         *int32
	frame        *int32
	frameCounter *int32
	speed        *int32
	controller   sdl.JoystickID
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
	h = h / 4
	p.H = &h
	w = w / 4
	p.W = &w
	p.size = p.H
	speed := int32(5)
	p.speed = &speed
	p.frame = new(int32)
	p.frameCounter = new(int32)
	p.dirx = new(int32)
	p.diry = new(int32)
	p.dir = new(int32)
	return p, nil
}

func mustPlayer(p *Player, err error) *Player {
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Player) Paint(r *sdl.Renderer) error {
	v := r.GetViewport()
	size := atomic.LoadInt32(p.size)
	return r.Copy(p.texture, p.getFrame(), &sdl.Rect{X: v.W/2 - atomic.LoadInt32(p.W)/2, Y: v.H/2 - atomic.LoadInt32(p.H)/2, W: size, H: size})
}

func (p *Player) Update() {
	framCounter := atomic.AddInt32(p.frameCounter, 1)

	dirx := atomic.LoadInt32(p.dirx) / 2000
	diry := atomic.LoadInt32(p.diry) / 2000
	fmt.Printf("%v \t%v\r", dirx, diry)
	if dirx == 0 && diry == 0 {
		return
	}
	var dir int32
	switch {
	case dirx == 3:
		dir = 0
	case dirx == -3:
		dir = 1
	case diry == 3:
		dir = 2
	case diry == -3:
		dir = 3
	}
	atomic.StoreInt32(p.dir, dir)

	if framCounter%atomic.LoadInt32(p.speed) == 0 {
		atomic.StoreInt32(p.frame, (atomic.LoadInt32(p.frame)+1)%4)
	}
	if framCounter > 10000 {
		atomic.StoreInt32(p.frameCounter, 0)
	}
}

func (p *Player) Handle(evt sdl.Event) bool {
	switch e := evt.(type) {
	case *sdl.ControllerButtonEvent:
	case *sdl.ControllerAxisEvent:
		if e.Which == p.controller {
			switch e.Axis {
			case sdl.CONTROLLER_AXIS_LEFTX, sdl.CONTROLLER_AXIS_RIGHTX:
				// fmt.Println("CONTROLLER_AXIS_LEFTX")
				atomic.StoreInt32(p.dirx, int32(e.Value))
			case sdl.CONTROLLER_AXIS_LEFTY, sdl.CONTROLLER_AXIS_RIGHTY:
				// fmt.Println("CONTROLLER_AXIS_LEFTY")
				atomic.StoreInt32(p.diry, int32(e.Value))
			// case sdl.CONTROLLER_AXIS_RIGHTX:
			// 	fmt.Println("CONTROLLER_AXIS_RIGHTX")
			// case sdl.CONTROLLER_AXIS_RIGHTY:
			// 	fmt.Println("CONTROLLER_AXIS_RIGHTY")
			case sdl.CONTROLLER_AXIS_MAX:
				fmt.Println("CONTROLLER_AXIS_MAX")
			case sdl.CONTROLLER_AXIS_TRIGGERLEFT:
				fmt.Println("CONTROLLER_AXIS_TRIGGERLEFT")
			case sdl.CONTROLLER_AXIS_TRIGGERRIGHT:
				fmt.Println("CONTROLLER_AXIS_TRIGGERRIGHT")
			default:
				fmt.Printf("Axis: %v Type: %v Val:%v\n", e.Axis, e.Type, e.Value)
			}
		}
		return true
	case *sdl.ControllerDeviceEvent:

	case *sdl.JoyAxisEvent:
	case *sdl.JoyButtonEvent:
	case *sdl.JoyDeviceEvent:

	}
	return false
}

func (p *Player) getFrame() *sdl.Rect {
	r := &sdl.Rect{
		W: atomic.LoadInt32(p.W),
		H: atomic.LoadInt32(p.H),
		X: atomic.LoadInt32(p.frame) * atomic.LoadInt32(p.W),
		Y: atomic.LoadInt32(p.dir) * atomic.LoadInt32(p.H),
	}
	return r
}

func (p *Player) Destroy() {

}

func (p *Player) Reset() {

}

func (p *Player) Resize(delta int32) {
	size := atomic.AddInt32(p.size, delta)
	if size > 256 {
		atomic.StoreInt32(p.size, 256)
	}
	if size < 10 {
		atomic.StoreInt32(p.size, 10)
	}
	fmt.Println("Player size:", size)
}
