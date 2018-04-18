package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	// FramesPerSecond num of frames a second to try and render
	FramesPerSecond = 60
	// FrameTolerance howmany frames per second are allowed to be skipped without recalculating frame time.
	FrameTolerance = FramesPerSecond / 10
)

func main() {
	time.AfterFunc(time.Second*60, func() {
		os.Exit(0)
	})

	r, cleanup := GetRenderer(800, 1000)
	defer cleanup()

	w := NewWorld(100, 100, 30, r)

	entities := EntityList{w}

	go func() {
		for range time.Tick(time.Second / 1) {
			// w.ShuffleTiles()
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

	if err := ttf.Init(); err != nil {
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

	// go eventLoop()
	return r,
		func() {
			window.Destroy()
			sdl.Quit()
			ttf.Quit()
			runtime.UnlockOSThread()
		}

}
