package main

import (
	"fmt"
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
	// return &p.cells
	p.RLock()
	if rect == p.prevView {
		p.RUnlock()
		return &p.viewBuffer
	}
	p.RUnlock()
	p.Lock()
	startX := p.originx + rect.X
	startY := p.originy + rect.Y
	rectCp := rect
	rectCp.X = startX
	rectCp.Y = startY
	if p.OutofBounds(rectCp) {
	}
	r := &rngSource{}
	r.Seed(24234)
	viewBuffer := make([][]Tile, rect.H)
	for i := int32(0); i < rect.H; i++ {
		// fmt.Printf("len:%v start:%v end:%v\n", len(view[row]), startX-w/2, startX+w/2)
		y := i + startY - (rect.H / 2)
		row := make([]Tile, rect.W)
		for j := int32(0); j < rect.W; j++ {
			x := j + startX - (rect.W / 2)
			row[j] = Tile{texture: r.GetInt(int(x), int(y), 256)}
		}
		viewBuffer[i] = row
	}
	p.prevView = rect
	p.viewBuffer = viewBuffer
	p.Unlock()
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

func NewPlane(textures []*TileRender, h int, w int) *Plane {
	p := new(Plane)
	p.originx = int32(w) / 2
	p.originy = int32(h) / 2
	p.cells = NewTiles(h, w)
	p.size = (h + w) / 2
	// p.ShuffleTiles(textures)
	p.cells = genTiles(p.size, 24325, textures)
	return p
}

func NewTiles(h, w int) [][]Tile {
	tiles := make([][]Tile, h)
	for row := range tiles {
		tiles[row] = make([]Tile, w)
	}
	return tiles
}
