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

type World struct {
	cells    [][]Tile
	tileSize int32
	textures []*sdl.Texture
	sync.RWMutex
}

type Tile struct {
	Name    string
	texture int
}

var _ Painter = &World{}

// Paint will draw the map in the renderer
func (w *World) Paint(r *sdl.Renderer) (err error) {
	w.RLock()
	defer w.RUnlock()

	for row, cols := range w.cells {
		for col, t := range cols {
			if t.texture > len(w.textures)-1 || w.textures[t.texture] == nil {
				fmt.Println("Texture not found:", t.texture)
				continue
			}
			err = r.Copy(w.textures[t.texture], nil, &sdl.Rect{H: w.tileSize, W: w.tileSize, X: int32(col) * w.tileSize, Y: int32(row) * w.tileSize})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewWorld(h, w int, tileSize int32, textures []*sdl.Texture) *World {
	world := new(World)
	world.cells = make([][]Tile, h)
	world.textures = textures
	world.tileSize = tileSize
	for row := range world.cells {
		world.cells[row] = make([]Tile, w)
	}
	world.ShuffleTiles()
	return world
}

func (w *World) ShuffleTiles() {
	w.Lock()
	defer w.Unlock()
	for row := range w.cells {
		for col := range w.cells[row] {
			w.cells[row][col].texture = int(rand.Int63()) % (len(w.textures) - 1)
		}
	}
}

func (w *World) ChangeTileSize(delta int32) {
	w.Lock()
	w.Unlock()
	w.tileSize += delta
	if w.tileSize > 256 {
		w.tileSize = 256
	}
	if w.tileSize < 10 {
		w.tileSize = 10
	}
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
