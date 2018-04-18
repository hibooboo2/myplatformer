package main

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type Plane struct {
	sync.RWMutex
	viewBuffer []Tile
	tiles      map[sdl.Point]Tile
	prevView   Rect
	rnd        *rngSource
}

func (p *Plane) View(x, y, h, w int32) *[]Tile {
	p.RLock()
	if x == p.prevView.X && p.prevView.Y == y && p.prevView.H == h && p.prevView.W == w {
		p.RUnlock()
		return &p.viewBuffer
	}
	p.RUnlock()

	p.Lock()
	tiles := []Tile{}
	for i := 0; i < int(h); i++ {
		yLoc := y - (h / 2) + int32(i)
		for j := 0; j < int(w); j++ {
			xLoc := x - (w / 2) + int32(j)
			t := Tile{texture: p.rnd.GetInt(int(xLoc), int(yLoc), 236), Loc: sdl.Point{X: int32(j), Y: int32(i)}}
			tiles = append(tiles, t)
		}
	}
	p.viewBuffer = tiles
	p.prevView = Rect{sdl.Rect{X: x, Y: y, W: w, H: h}}
	p.Unlock()

	return &p.viewBuffer
}

func NewPlane(textures []*TileRender, h int, w int) *Plane {
	p := new(Plane)
	p.rnd = new(rngSource)
	p.rnd.Seed(234524356)
	return p
}

func NewTiles(h, w int) [][]Tile {
	tiles := make([][]Tile, h)
	for row := range tiles {
		tiles[row] = make([]Tile, w)
	}
	return tiles
}
