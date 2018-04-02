package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	events = make(chan sdl.Event, 2)
	quit   = make(chan struct{}, 2)
	//Handle list of event Handlers the engine uses.
	handlers       []EventHandler
	handlersLock   = sync.RWMutex{}
	DefaultHandler = EventHandlerFunc(handleEvent)
)

func EventLoop(r *sdl.Renderer, e Painter) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	go func(e Painter) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		start := time.Now()
		i := 0
		computedFrameTime := time.Second / time.Duration(FramesPerSecond)
		for {
			if computedFrameTime <= 0 {
				computedFrameTime = time.Nanosecond
			}
			ticker := time.NewTicker(computedFrameTime)
			for range ticker.C {
				e.Paint(r)
				r.Present()
				i++
				// fmt.Print(".")
				took := time.Since(start)
				if took > time.Second {
					avg := took / time.Duration(i)
					adjust := false
					buff := float64(took / time.Second)
					expFrames := int(float64(FramesPerSecond) * buff)
					if expFrames > i+int(buff*FrameTolerance) || expFrames < i-int(buff*FrameTolerance) {
						fmt.Printf("Missed frames ExpFrames: %v Actual: %v\n", expFrames, i)
						adjust = true
					} else {
						fmt.Printf("no missed frames Exp: %v Got: %v\r", expFrames, i)
					}
					if adjust {
						ticker.Stop()
						computedFrameTime = computedFrameTime - (avg - time.Second/FramesPerSecond)
						fmt.Println(computedFrameTime)
					}
					start = time.Now()
					i = 0
					break
				}

			}
		}
	}(e)

	for {
		select {
		case event := <-events:
			if event == nil {
				continue
			}
			handled := false
			handlersLock.RLock()
			defer handlersLock.RUnlock()
			for _, h := range handlers {
				handled = h.Handle(event)
				if handled {
					break
				}
			}
			if !handled {
				handled = DefaultHandler.Handle(event)
			}
			if !handled {
				fmt.Printf("Unhandled %T %##v\n", event, event)
			}
		case events <- sdl.PollEvent():
		case <-quit:
			fmt.Println("Quitting")
			return
		}
	}
}

type EventHandler interface {
	Handle(evt sdl.Event) bool
}

type EventHandlerFunc func(evt sdl.Event) bool

func (eh EventHandlerFunc) Handle(evt sdl.Event) bool {
	return eh(evt)
}

func handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case nil:
	case *sdl.QuitEvent:
		quit <- struct{}{}
		return true
	case *sdl.TextInputEvent:
	case *sdl.KeyboardEvent:
		if (e.Keysym.Mod | sdl.KMOD_CTRL) == sdl.KMOD_CTRL {
			switch e.Keysym.Sym {
			case sdl.K_w, sdl.K_q, sdl.K_c, sdl.K_d:
				quit <- struct{}{}
				return true
			}
		}
	case *sdl.MouseMotionEvent:
	case *sdl.WindowEvent:
	case *sdl.MouseButtonEvent:
	default:
		fmt.Printf("%T\n", event)
		return false
	}
	return true
}

func AddHandler(h EventHandler) {
	handlersLock.Lock()
	defer handlersLock.Unlock()
	handlers = append(handlers, h)
}

func AddHandlerFunc(h func(evt sdl.Event) bool) {
	AddHandler(EventHandlerFunc(h))
}
