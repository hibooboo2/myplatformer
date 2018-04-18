package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	_ Entity  = &World{}
	_ Painter = &World{}
)

type World struct {
	cells    *Plane
	players  map[sdl.JoystickID]*Player
	tileSize *int32
	textures []*TileRender
	sync.RWMutex
	paused bool
	r      *sdl.Renderer
}

// Paint will draw the map in the renderer
func (w *World) Paint(r *sdl.Renderer) (err error) {
	w.RLock()
	defer w.RUnlock()
	v := r.GetViewport()
	tileSize := atomic.LoadInt32(w.tileSize)
	tilesWide := v.W / tileSize
	tilesHigh := v.H / tileSize

	view := w.cells.View(atomic.LoadInt32(w.players[999].X)/tileSize,
		atomic.LoadInt32(w.players[999].Y)/tileSize,
		tilesHigh*2,
		tilesWide*2,
	)

	mainTexture := r.GetRenderTarget()

	var tileText *sdl.Texture
	tileText, err = r.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_TARGET, v.W, v.H)
	if err != nil {
		return err
	}
	defer tileText.Destroy()
	// tileText.SetBlendMode(sdl.BLENDMODE_BLEND)
	// tileText.SetAlphaMod(100)

	err = r.SetRenderTarget(tileText)
	if err != nil {
		return err
	}

	// tileText.SetColorMod(100, 100, 100)

	for _, t := range *view {
		if t.texture > len(w.textures)-1 || w.textures[t.texture] == nil {
			fmt.Println("Texture not found:", t.texture)
			continue
		}

		tile := &sdl.Rect{H: tileSize, W: tileSize, X: t.Loc.X * tileSize, Y: t.Loc.Y * tileSize}

		err = r.Copy(w.textures[t.texture].texture, nil, tile)
		if err != nil {
			return err
		}

	}

	if err = r.SetRenderTarget(mainTexture); err != nil {
		return err
	}

	err = r.Copy(tileText, nil, &v)
	if err != nil {
		return err
	}

	for _, p := range w.players {
		err = p.Paint(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *World) Resize(delta int32) {
	tileSize := atomic.LoadInt32(w.tileSize) + delta
	if tileSize > 256 {
		tileSize = 256
	}
	if tileSize < 10 {
		tileSize = 10
	}
	atomic.StoreInt32(w.tileSize, tileSize)
	fmt.Println("World Tile Size:", tileSize)
}

func (w *World) Destroy() {
}

func (w *World) Reset() {
}

func (w *World) Update() {
	w.RLock()
	if w.paused {
		w.RUnlock()
		return
	}
	w.RUnlock()
	for _, p := range w.players {
		p.Update()
	}
}

func (w *World) Handle(evt sdl.Event) bool {
	var handled bool
	for _, p := range w.players {
		if p.Handle(evt) {
			handled = true
		}
	}
	switch e := evt.(type) {
	case *sdl.MouseWheelEvent:
		w.Resize(e.Y)
		return true
	case *sdl.ControllerDeviceEvent:
		switch e.Type {
		case sdl.CONTROLLERDEVICEADDED:
			w.Lock()
			p := mustPlayer(NewPlayer(w.r, sdl.GameControllerOpen(int(e.Which)), nil))
			w.players[p.controller] = p
			w.Unlock()
		case sdl.CONTROLLERDEVICEREMOVED:
			w.Lock()
			delete(w.players, e.Which)
			w.Unlock()
		}
	}
	return handled
}

func NewWorld(w int, h int, tileSize int32, r *sdl.Renderer) *World {
	world := new(World)
	world.textures = mustTexture(getTextures(r))
	world.tileSize = &tileSize
	world.cells = NewPlane(world.textures, h, w)
	world.r = r
	world.players = make(map[sdl.JoystickID]*Player)
	world.players[999] = mustPlayer(NewPlayer(r, nil, WSAD_KEYS))
	world.players[999].main = true
	world.players[998] = mustPlayer(NewPlayer(r, nil, ARROW_KEYS))
	return world
}
