package main

import (
	"fmt"
	"sync/atomic"

	. "github.com/hibooboo2/myplatformer/dir"
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
	frame        *int32
	frameCounter *int32
	speed        *int32
	controller   sdl.JoystickID
	keys         []uint8
}

func NewPlayer(r *sdl.Renderer) (*Player, error) {
	p := new(Player)

	t, err := img.LoadTexture(r, "assets/penguin.png")
	if err != nil {
		return nil, err
	}
	p.texture = t
	_, _, w, h, err := t.Query()
	if err != nil {
		return nil, err
	}
	h = h / 8
	p.H = &h
	w = w / 8
	p.W = &w
	p.size = p.H
	speed := int32(5)
	p.speed = &speed
	p.frame = new(int32)
	p.frameCounter = new(int32)
	p.dir = new(int32)
	*p.dir = STOP
	p.keys = sdl.GetKeyboardState()
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

	if atomic.LoadInt32(p.dir) == 99 {
		return
	}
	if framCounter%atomic.LoadInt32(p.speed) == 0 {
		atomic.StoreInt32(p.frame, (atomic.LoadInt32(p.frame)+1)%8)
	}
	if framCounter > 10000 {
		atomic.StoreInt32(p.frameCounter, 0)
	}
}

var DIR_FRAME = map[DIRECTION]int32{
	SOUTH:      0,
	SOUTH_WEST: 1,
	WEST:       2,
	NORTH_WEST: 3,
	NORTH:      4,
	NORTH_EAST: 5,
	EAST:       6,
	SOUTH_EAST: 7,
	STOP:       8,
}

func (p *Player) Handle(evt sdl.Event) bool {
	switch e := evt.(type) {
	case *sdl.ControllerButtonEvent:
	case *sdl.MouseWheelEvent:
		p.Resize(e.Y)
		return true
	case *sdl.ControllerAxisEvent:
		if e.Which == p.controller {
			switch e.Axis {
			case sdl.CONTROLLER_AXIS_LEFTX, sdl.CONTROLLER_AXIS_RIGHTX:
				// fmt.Println("CONTROLLER_AXIS_LEFTX")
				dir := DIRECTION(atomic.LoadInt32(p.dir))
				val := e.Value / AXIS_TOLERANCE
				switch {
				case val > 0:
					dir = dir & (NORTH | SOUTH)
					dir = dir | EAST
				case val < 0:
					dir = dir & (NORTH | SOUTH)
					dir = dir | WEST
				case val == 0:
					dir = STOP
				}
				fmt.Printf("%s %v X\n", DIRECTION(dir), val)
				atomic.StoreInt32(p.dir, int32(DIR_FRAME[dir]))
			case sdl.CONTROLLER_AXIS_LEFTY, sdl.CONTROLLER_AXIS_RIGHTY:
				// fmt.Println("CONTROLLER_AXIS_LEFTY")
				dir := DIRECTION(atomic.LoadInt32(p.dir))
				val := e.Value / AXIS_TOLERANCE
				switch {
				case val > 0:
					dir = dir & (EAST | WEST)
					dir = dir | NORTH
				case val < 0:
					dir = dir & (EAST | WEST)
					dir = dir | SOUTH
				case val == 0:
					dir = STOP
				}
				fmt.Printf("%s %v Y\n", DIRECTION(dir), val)
				atomic.StoreInt32(p.dir, int32(DIR_FRAME[dir]))
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
	case *sdl.KeyboardEvent:
		var dir DIRECTION
		switch {
		case (p.keys[sdl.SCANCODE_W] & p.keys[sdl.SCANCODE_A]) == 1:
			dir = NORTH_EAST
		case (p.keys[sdl.SCANCODE_W] & p.keys[sdl.SCANCODE_D]) == 1:
			dir = NORTH_WEST
		case (p.keys[sdl.SCANCODE_S] & p.keys[sdl.SCANCODE_A]) == 1:
			dir = SOUTH_EAST
		case (p.keys[sdl.SCANCODE_S] & p.keys[sdl.SCANCODE_D]) == 1:
			dir = SOUTH_WEST
		case p.keys[sdl.SCANCODE_W] == 1:
			dir = NORTH
		case p.keys[sdl.SCANCODE_S] == 1:
			dir = SOUTH
		case p.keys[sdl.SCANCODE_D] == 1:
			dir = WEST
		case p.keys[sdl.SCANCODE_A] == 1:
			dir = EAST
		default:
			dir = STOP
		}
		atomic.StoreInt32(p.dir, int32(DIR_FRAME[dir]))
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
