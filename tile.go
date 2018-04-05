package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	_ Entity = &World{}
)

type World struct {
	cells    *Plane
	players  map[sdl.JoystickID]*Player
	tileSize *int32
	height   int
	width    int
	textures []*sdl.Texture
	sync.RWMutex
	paused bool
	r      *sdl.Renderer
}

type Tile struct {
	Name    string
	texture int
	Loc     sdl.Point
}

var _ Painter = &World{}

// Paint will draw the map in the renderer
func (w *World) Paint(r *sdl.Renderer) (err error) {
	w.RLock()
	v := r.GetViewport()
	tileSize := atomic.LoadInt32(w.tileSize)
	tilesWide := v.W / tileSize
	tilesHigh := v.H / tileSize

	view := w.cells.View(sdl.Rect{
		X: atomic.LoadInt32(w.players[999].X) / tileSize,
		Y: atomic.LoadInt32(w.players[999].Y) / tileSize,
		H: tilesHigh + 1,
		W: tilesWide + 1,
	})

	for row, cols := range *view {
		for col, t := range cols {
			if t.texture > len(w.textures)-1 || w.textures[t.texture] == nil {
				fmt.Println("Texture not found:", t.texture)
				continue
			}
			tile := &sdl.Rect{H: tileSize, W: tileSize, X: int32(col) * tileSize, Y: int32(row) * tileSize}
			err = r.Copy(w.textures[t.texture], nil, tile)
			if err != nil {
				w.RUnlock()
				return err
			}
			// locT := renderText(r, fmt.Sprintf("%d,%d", t.Loc.X, t.Loc.Y))
			// err = r.Copy(locT, nil, &sdl.Rect{
			// 	H: tileSize / 2,
			// 	W: tileSize / 2,
			// })
			// if err != nil {
			// 	w.RUnlock()
			// 	return err
			// }
		}
	}
	if err := r.SetRenderTarget(nil); err != nil {
		return err
	}
	for _, p := range w.players {
		err = p.Paint(r)
		if err != nil {
			return err
		}
	}

	w.RUnlock()
	return nil
}

var (
	textOnce    sync.Once
	locTextures = map[string]*sdl.Texture{}
	font        *ttf.Font
)

func renderText(r *sdl.Renderer, text string) *sdl.Texture {
	r.Clear()
	textOnce.Do(func() {
		f, err := ttf.OpenFont("assets/fonts/Flappy.ttf", 20)
		if err != nil {
			fmt.Printf("could not load font: %v\n", err)
			return
		}
		font = f
	})
	t := locTextures[text]
	if t != nil {
		return t
	}

	c := sdl.Color{R: 200, G: 100, B: 100, A: 128}
	s, err := font.RenderUTF8Blended(text, c)
	if err != nil {
		fmt.Printf("could not render title: %v\n", err)
		return nil
	}
	// defer s.Free()

	t, err = r.CreateTextureFromSurface(s)
	if err != nil {
		fmt.Printf("could not create texture: %v\n", err)
		return nil
	}

	t.SetBlendMode(sdl.BLENDMODE_BLEND)

	locTextures[text] = t

	return t
}

func NewWorld(w int, h int, tileSize int32, r *sdl.Renderer) *World {
	world := new(World)
	world.textures = mustTexture(getTextures(r))
	world.tileSize = &tileSize
	world.height = h
	world.width = w
	world.cells = NewPlane(world.textures)
	world.r = r
	world.players = make(map[sdl.JoystickID]*Player)
	world.players[999] = mustPlayer(NewPlayer(r, nil, WSAD_KEYS))
	world.players[998] = mustPlayer(NewPlayer(r, nil, ARROW_KEYS))
	return world
}

func (w *World) newTileSlice() [][]Tile {
	w.RLock()
	height := w.height
	width := w.width
	w.RUnlock()

	tiles := make([][]Tile, height)
	for row := range tiles {
		tiles[row] = make([]Tile, width)
	}
	return tiles
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

func mustTexture(textures []*sdl.Texture, err error) []*sdl.Texture {
	if err != nil {
		panic(err)
	}
	return textures
}

func getTextures(r *sdl.Renderer) ([]*sdl.Texture, error) {
	var textures []*sdl.Texture
	dir := "assets/images/lostgarden/1/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".BMP") {
			continue
		}

		t, err := img.LoadTexture(r, dir+f.Name())
		if err != nil {
			return nil, err
		}
		t.SetBlendMode(sdl.BLENDMODE_BLEND)
		textures = append(textures, t)
	}
	return textures, nil
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
