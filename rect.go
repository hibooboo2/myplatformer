package main

import "github.com/veandco/go-sdl2/sdl"

type Rect struct {
	sdl.Rect
}

func (r *Rect) Contains(x, y int32) bool {
	return (x > r.X-r.W/2) &&
		(x < r.X+r.W/2) &&
		(y > r.Y-r.H/2) &&
		(y < r.Y+r.H/2)
}
