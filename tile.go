package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	_ Entity = &World{}
)

type World struct {
	cells    [][]Tile
	players  []*Player
	tileSize int32
	height   int
	width    int
	textures []*sdl.Texture
	sync.RWMutex
	paused bool
}
type Tile struct {
	Name    string
	texture int
}

var _ Painter = &World{}

// Paint will draw the map in the renderer
func (w *World) Paint(r *sdl.Renderer) (err error) {
	w.RLock()
	v := r.GetViewport()
	tilesWide := v.W / w.tileSize
	tilesHigh := v.H / w.tileSize

	for row, cols := range w.cells[:tilesHigh+1] {
		for col, t := range cols[:tilesWide+1] {
			if t.texture > len(w.textures)-1 || w.textures[t.texture] == nil {
				fmt.Println("Texture not found:", t.texture)
				continue
			}
			err = r.Copy(w.textures[t.texture], nil, &sdl.Rect{H: w.tileSize, W: w.tileSize, X: int32(col) * w.tileSize, Y: int32(row) * w.tileSize})
			if err != nil {
				w.RUnlock()
				return err
			}
		}
	}
	w.RUnlock()
	return nil
}

func NewWorld(w int, h int, tileSize int32, r *sdl.Renderer) *World {
	world := new(World)
	world.textures = mustTexture(getTextures(r))
	world.tileSize = tileSize
	world.height = h
	world.width = w
	world.cells = world.newTileSlice()
	world.ShuffleTiles()
	p := mustPlayer(NewPlayer(r))
	world.players = append(world.players, p)
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

func (w *World) ShuffleTiles() {
	newTiles := w.newTileSlice()

	for row := range newTiles {
		for col := range newTiles[row] {
			newTiles[row][col].texture = int(rand.Int63()) % (len(w.textures) - 1)
		}
	}

	w.Lock()
	w.cells = newTiles
	w.Unlock()
}

func (w *World) Resize(delta int32) {
	w.Lock()
	defer w.Unlock()
	w.tileSize += delta
	if w.tileSize > 256 {
		w.tileSize = 256
	}
	if w.tileSize < 10 {
		w.tileSize = 10
	}
	fmt.Println("World Tile Size:", w.tileSize)
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
		textures = append(textures, t)
	}
	return textures, nil
}

func (w *World) Destroy() {
}

func (w *World) Reset() {
}

func (w *World) Update() {
}

func (w *World) Handle(evt sdl.Event) bool {
	return false
}
