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
	X, Y         *int32
	H, W         *int32
	size         *int32
	dir          *int32
	dirx         *int32
	diry         *int32
	frame        *int32
	moving       *int32
	frameCounter *int32
	speed        *int32
	controller   sdl.JoystickID
	keys         []uint8
	con          *sdl.GameController
	keyControls  map[DIRECTION]sdl.Keycode
}

func NewPlayer(r *sdl.Renderer, c *sdl.GameController, controls map[DIRECTION]sdl.Keycode) (*Player, error) {
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
	p.H = new(int32)
	p.W = new(int32)
	*p.H = h / 8
	*p.W = w / 8
	p.size = new(int32)
	*p.size = *p.H
	speed := int32(5)
	p.speed = &speed
	p.frame = new(int32)
	p.frameCounter = new(int32)
	p.dir = new(int32)
	*p.dir = int32(STOP)
	p.dirx = new(int32)
	p.diry = new(int32)
	p.X = new(int32)
	p.Y = new(int32)
	p.moving = new(int32)
	p.keys = sdl.GetKeyboardState()
	p.controller = c.Joystick().InstanceID()
	p.con = c
	p.keyControls = controls
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
	return r.Copy(p.texture, p.getFrame(), &sdl.Rect{
		X: v.W/2 - atomic.LoadInt32(p.size)/2 + atomic.LoadInt32(p.X),
		Y: v.H/2 - atomic.LoadInt32(p.size)/2 + atomic.LoadInt32(p.Y),
		W: size,
		H: size,
	})
}

func (p *Player) Dir() DIRECTION {
	return frameToDirPT(p.dir)
}

func frameToDir(frame int32) DIRECTION {
	for d, f := range DIR_FRAME {
		if f == frame {
			return d
		}
	}
	return STOP
}

func frameToDirPT(framePt *int32) DIRECTION {
	return frameToDir(atomic.LoadInt32(framePt))
}

func (p *Player) DirVals() {
	// x := atomic.LoadInt32(p.dirx)
	// y := atomic.LoadInt32(p.diry)
	// fmt.Printf("%d..%d..\t\t%s..............\r", x, y, p.Dir().String())
}

func (p *Player) ChkKey() {
	if p.keyControls == nil {
		return
	}

	var dir DIRECTION
	switch {
	case (p.keys[p.keyControls[NORTH]] & p.keys[p.keyControls[EAST]]) == 1:
		dir = NORTH_EAST
	case (p.keys[p.keyControls[NORTH]] & p.keys[p.keyControls[WEST]]) == 1:
		dir = NORTH_WEST
	case (p.keys[p.keyControls[SOUTH]] & p.keys[p.keyControls[EAST]]) == 1:
		dir = SOUTH_EAST
	case (p.keys[p.keyControls[SOUTH]] & p.keys[p.keyControls[WEST]]) == 1:
		dir = SOUTH_WEST
	case p.keys[p.keyControls[NORTH]] == 1:
		dir = NORTH
	case p.keys[p.keyControls[SOUTH]] == 1:
		dir = SOUTH
	case p.keys[p.keyControls[WEST]] == 1:
		dir = WEST
	case p.keys[p.keyControls[EAST]] == 1:
		dir = EAST
	default:
		dir = STOP
	}
	p.SetDir(dir)
	p.DirVals()
}

// SetDir Sets the direction of the character
func (p *Player) SetDir(dir DIRECTION) {
	if dir == STOP {
		atomic.StoreInt32(p.moving, int32(STOP))
	} else {
		atomic.StoreInt32(p.moving, 0)
		atomic.StoreInt32(p.dir, int32(DIR_FRAME[dir]))
	}
}

func (p *Player) ChkCon() {
	if p.con == nil {
		fmt.Println("No con")
		return
	}
	dirx := atomic.LoadInt32(p.dirx)
	diry := atomic.LoadInt32(p.diry)
	var dir DIRECTION
	switch {
	case dirx > AXIS_TOLERANCE && diry > AXIS_TOLERANCE:
		dir = NORTH_EAST
	case dirx > AXIS_TOLERANCE && diry < -AXIS_TOLERANCE:
		dir = SOUTH_EAST
	case dirx < -AXIS_TOLERANCE && diry > AXIS_TOLERANCE:
		dir = NORTH_WEST
	case dirx < -AXIS_TOLERANCE && diry < -AXIS_TOLERANCE:
		dir = SOUTH_WEST
	case dirx > AXIS_TOLERANCE:
		dir = EAST
	case dirx < -AXIS_TOLERANCE:
		dir = WEST
	case diry > AXIS_TOLERANCE:
		dir = NORTH
	case diry < -AXIS_TOLERANCE:
		dir = SOUTH
	default:
		dir = STOP
	}
	p.SetDir(dir)
	p.DirVals()
}

func (p *Player) Update() {
	framCounter := atomic.AddInt32(p.frameCounter, 1)

	if DIRECTION(atomic.LoadInt32(p.moving)) == STOP {
		return
	}

	spd := atomic.LoadInt32(p.speed)
	switch p.Dir() {
	case NORTH:
		atomic.AddInt32(p.Y, -spd)
	case NORTH_EAST:
		atomic.AddInt32(p.Y, -spd)
		atomic.AddInt32(p.X, -spd)
	case EAST:
		atomic.AddInt32(p.X, -spd)
	case SOUTH_EAST:
		atomic.AddInt32(p.X, -spd)
		atomic.AddInt32(p.Y, spd)
	case SOUTH:
		atomic.AddInt32(p.Y, spd)
	case SOUTH_WEST:
		atomic.AddInt32(p.Y, spd)
		atomic.AddInt32(p.X, spd)
	case WEST:
		atomic.AddInt32(p.X, spd)
	case NORTH_WEST:
		atomic.AddInt32(p.Y, -spd)
		atomic.AddInt32(p.X, spd)
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

var (
	WSAD_KEYS = map[DIRECTION]sdl.Keycode{
		NORTH: sdl.SCANCODE_W,
		SOUTH: sdl.SCANCODE_S,
		EAST:  sdl.SCANCODE_A,
		WEST:  sdl.SCANCODE_D,
	}
	ARROW_KEYS = map[DIRECTION]sdl.Keycode{
		NORTH: sdl.SCANCODE_UP,
		SOUTH: sdl.SCANCODE_DOWN,
		EAST:  sdl.SCANCODE_LEFT,
		WEST:  sdl.SCANCODE_RIGHT,
	}
)

func (p *Player) Handle(evt sdl.Event) bool {
	switch e := evt.(type) {
	case *sdl.ControllerButtonEvent:
	case *sdl.MouseWheelEvent:
		p.Resize(e.Y)
		return true
	case *sdl.ControllerAxisEvent:
		if p.con != nil && e.Which == p.controller {
			switch e.Axis {
			case sdl.CONTROLLER_AXIS_LEFTX, sdl.CONTROLLER_AXIS_RIGHTX:
				atomic.StoreInt32(p.dirx, -int32(e.Value))
				p.ChkCon()
			case sdl.CONTROLLER_AXIS_LEFTY, sdl.CONTROLLER_AXIS_RIGHTY:
				atomic.StoreInt32(p.diry, -int32(e.Value))
				p.ChkCon()
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
		// fmt.Printf("keySym:%v sym:%v\n", e.Keysym, e.Keysym.Sym)
		p.ChkKey()
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
}
