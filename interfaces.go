package main

import "github.com/veandco/go-sdl2/sdl"

type Destroyer interface {
	Destroy()
}

type Painter interface {
	Paint(r *sdl.Renderer) error
}

type Updater interface {
	Update()
}

type Collider interface {
	Collides(c Collider) bool
	Type() string
	CollidingTypes() []string
}

type Reseter interface {
	Reset()
}

type Entity interface {
	Destroyer
	Painter
	Updater
	Reseter
}

type EntityList []Entity

var _ Entity = &EntityList{}

func (el *EntityList) Destroy() {
	for _, e := range *el {
		e.Destroy()
	}
}

func (el *EntityList) Paint(r *sdl.Renderer) error {
	for _, e := range *el {
		err := e.Paint(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (el *EntityList) Update() {
	for _, e := range *el {
		e.Update()
	}
}

func (el *EntityList) Reset() {
	for _, e := range *el {
		e.Reset()
	}
}
