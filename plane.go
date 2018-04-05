package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type Plane struct {
	sync.RWMutex
	cells            [][]Tile
	originx, originy int32
	size             int
	viewBuffer       [][]Tile
	prevView         sdl.Rect
}

func (p *Plane) View(rect sdl.Rect) *[][]Tile {
	p.RLock()
	if rect != p.prevView {
		startX := p.originx + rect.X
		startY := p.originy + rect.Y
		rectCp := rect
		rectCp.X = startX
		rectCp.Y = startY
		if p.OutofBounds(rectCp) {
		}
		viewBuffer := make([][]Tile, rect.H)
		for i := int32(0); i < rect.H; i++ {
			// fmt.Printf("len:%v start:%v end:%v\n", len(view[row]), startX-w/2, startX+w/2)
			viewBuffer[i] = p.cells[i+startY][startX-rect.W/2 : startX+rect.W/2]
		}
		p.prevView = rect
		p.viewBuffer = viewBuffer
	}
	p.RUnlock()
	return &p.viewBuffer
}

func (p *Plane) OutofBounds(rect sdl.Rect) bool {
	fullView := &sdl.Rect{
		W: int32(len(p.cells[0])),
		H: int32(len(p.cells)),
		X: p.originx,
		Y: p.originy,
	}
	u := rect.Union(fullView)
	fmt.Println(rect, fullView, u)
	return false
}

func NewPlane(textures []*sdl.Texture) *Plane {
	size := int(200)
	p := new(Plane)
	p.originx = int32(size) / 2
	p.originy = int32(size) / 2
	p.size = size
	p.cells = NewTiles(size, size)
	p.ShuffleTiles(textures)
	return p
}

func NewTiles(h, w int) [][]Tile {
	tiles := make([][]Tile, h)
	for row := range tiles {
		tiles[row] = make([]Tile, w)
	}
	return tiles
}

func (p *Plane) ShuffleTiles(textures []*sdl.Texture) {
	newTiles := NewTiles(len(p.cells), len(p.cells[0]))

	for row := range newTiles {
		for col := range newTiles[row] {
			newTiles[row][col].Loc = sdl.Point{X: int32(col) - int32(p.size/2), Y: int32(row) - int32(p.size/2)}
			newTiles[row][col].texture = int(rand.Int63()) % (len(textures) - 1)
		}
	}

	p.Lock()
	p.cells = newTiles
	p.Unlock()
}
