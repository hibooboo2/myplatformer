package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// FramesPerSecond num of frames a second to try and render
	FramesPerSecond = 60
	// FrameTolerance howmany frames per second are allowed to be skipped without recalculating frame time.
	FrameTolerance = FramesPerSecond / 10
)

func main() {
	r, cleanup := GetRenderer(800, 600)
	defer cleanup()

	w := NewWorld(2000, 2000, 100, mustTexture(getTextures(r)))
	p, err := NewPlayer(r)
	if err != nil {
		panic(err)
	}
	entities := EntityList{w, p}
	go func() {
		for range time.Tick(time.Second / 1) {
			w.ShuffleTiles()
		}
	}()

	AddHandlerFunc(func(evt sdl.Event) bool {
		switch e := evt.(type) {
		case *sdl.JoyDeviceEvent:
			return true
		case *sdl.ControllerDeviceEvent:
			switch e.Type {
			case sdl.CONTROLLERDEVICEADDED:
				sdl.GameControllerOpen(int(e.Which))
			case sdl.CONTROLLERDEVICEREMOVED:
				fmt.Println("Removed controller")
			}
			return true
		case *sdl.MouseWheelEvent:
			entities.Resize(e.Y)
			return true
		case *sdl.ControllerButtonEvent:
			switch e.Type {
			case sdl.CONTROLLERBUTTONDOWN:
				fmt.Println("Button pressed down")
			case sdl.CONTROLLERBUTTONUP:
				fmt.Println("Button released")
			default:
				fmt.Printf("Unknown btn event: %##v\n", e)
			}
			return true
		case *sdl.ControllerAxisEvent:
			fmt.Printf("Ax:%v Type:%v Val:%v Controller:%v\n", e.Axis, e.Type, e.Value, e.Which)
			return true
		case *sdl.JoyAxisEvent, *sdl.JoyButtonEvent:
			return true
		default:
			return false
		}
	})

	EventLoop(r, &entities)
}

func GetRenderer(h, w int32) (*sdl.Renderer, func()) {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, r, err := sdl.CreateWindowAndRenderer(w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	window.SetResizable(true)
	window.SetBordered(true)
	// window.SetGrab(true)
	// window.SetWindowOpacity(0.4)
	sdl.GameControllerOpen(0)
	// go eventLoop()
	return r,
		func() {
			window.Destroy()
			sdl.Quit()
			runtime.UnlockOSThread()
		}

}
