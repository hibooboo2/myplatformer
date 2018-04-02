package main

import (
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// FrameTime time between Frames. Used for the paint Loop
	FramesPerSecond = 60
	FrameTolerance  = FramesPerSecond / 10
)

func main() {
	r, cleanup := Init(800, 600)
	defer cleanup()

	w := NewWorld(2000, 2000, 100, mustTexture(getTextures(r)))
	go func() {
		for range time.Tick(time.Second / 60) {
			w.ShuffleTiles()
		}
	}()

	// AddHandlerFunc(func(evt sdl.Event) bool {
	// 	switch e := evt.(type) {
	// 	case *sdl.MouseWheelEvent:
	// 		w.ChangeTileSize(e.Y)
	// 		fmt.Println(w.tileSize, e.Y)
	// 		return true
	// 	}
	// 	return false
	// })

	EventLoop(r, w)
}

func Init(h, w int32) (*sdl.Renderer, func()) {
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

	// go eventLoop()
	return r,
		func() {
			window.Destroy()
			sdl.Quit()
			runtime.UnlockOSThread()
		}

}
