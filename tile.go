package main

import (
	"fmt"
	"io/ioutil"
	"sort"
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
	textures []*TileRender
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
	defer w.RUnlock()
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

	mainTexture := r.GetRenderTarget()

	var tileText *sdl.Texture
	tileText, err = r.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_TARGET, v.W, v.H)
	if err != nil {
		return err
	}
	defer tileText.Destroy()
	tileText.SetBlendMode(sdl.BLENDMODE_BLEND)
	tileText.SetAlphaMod(100)

	err = r.SetRenderTarget(tileText)
	if err != nil {
		return err
	}

	tileText.SetColorMod(100, 100, 100)
	// r.SetDrawColor(0, 0, 0, 0)
	// r.FillRect(&v)

	for row, cols := range *view {
		for col, t := range cols {
			if t.texture > len(w.textures)-1 || w.textures[t.texture] == nil {
				// fmt.Println("Texture not found:", t.texture)
				continue
			}

			tile := &sdl.Rect{H: tileSize, W: tileSize, X: int32(col) * tileSize, Y: int32(row) * tileSize}

			// locT := renderText(r, fmt.Sprintf("%d,%d", t.Loc.X, t.Loc.Y))
			// err = r.Copy(locT, nil, tile)
			// if err != nil {
			// 	return err
			// }

			err = r.Copy(w.textures[t.texture].texture, nil, tile)
			if err != nil {
				return err
			}

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
	t.SetAlphaMod(255)

	locTextures[text] = t

	return t
}

func NewWorld(w int, h int, tileSize int32, r *sdl.Renderer) *World {
	world := new(World)
	world.textures = mustTexture(getTextures(r))
	world.tileSize = &tileSize
	world.height = h
	world.width = w
	world.cells = NewPlane(world.textures, h, w)
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

func mustTexture(textures []*TileRender, err error) []*TileRender {
	if err != nil {
		panic(err)
	}
	return textures
}

type TileRender struct {
	texture      *sdl.Texture
	Num          int
	Name         string
	Transition   string
	Group        string
	TileType     string
	ColorDepth   string
	TileVersion  string
	RandomFactor string
}

func getTextures(r *sdl.Renderer) ([]*TileRender, error) {
	var textures []*TileRender
	dir := "assets/images/lostgarden/1/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for i, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".BMP") {
			continue
		}
		textureName := dir + f.Name()
		tr, err := NewTileRender(r, textureName)
		if err == ErrIgnore {
			continue
		}
		tr.Num = i

		if err != nil {
			return nil, err
		}
		textures = append(textures, tr)
	}
	return textures, nil
}

var ErrIgnore = fmt.Errorf("Ignore Texture")

func NewTileRender(r *sdl.Renderer, textureName string) (*TileRender, error) {
	tr := new(TileRender)
	if r != nil {
		t, err := img.LoadTexture(r, textureName)
		if err != nil {
			return nil, err
		}
		t.SetBlendMode(sdl.BLENDMODE_BLEND)
		t.SetAlphaMod(200)
		tr.texture = t

	}

	tr.SetProps(textureName)

	switch tr.TileType {
	case "n", "s":
		return nil, ErrIgnore
	}
	// log.Printf("%v %v", tr.Group, tr.Transition)

	return tr, nil
}

// SetProps sets the string properties that are all derived from the name of the testure.
func (t *TileRender) SetProps(textureFileName string) {
	spl := strings.Split(textureFileName, "/")
	t.Name = strings.TrimSuffix(spl[len(spl)-1], ".BMP")
	if len(t.Name) != 8 {
		panic("Invalid Name For Texture" + textureFileName)
	}
	t.Group = t.Name[:2]
	t.Transition = t.Name[2:4]
	t.TileType = t.Name[4:5]
	t.ColorDepth = t.Name[5:6]
	t.TileVersion = t.Name[6:7]
	t.RandomFactor = t.Name[7:]
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

func genTiles(size, seed int, tilesCanUse []*TileRender) [][]Tile {
	//Pick fill or not for all groups random throughout cells
	//Place fill shapes
	//Transition all fils to adajacent types
	r := &rngSource{}
	r.Seed(seed)
	groups := usableGroups(tilesCanUse)
	groupKeys := []string{}
	for k := range groups {
		groupKeys = append(groupKeys, k)
	}
	sort.Slice(groupKeys, func(i int, j int) bool {
		return groupKeys[i] > groupKeys[j]
	})

	tiles := make([][]Tile, size)
	for i := range tiles {
		tileRow := make([]Tile, size)
		for j := range tileRow {
			group := groupKeys[r.GetInt(j, i, len(groupKeys))]
			tile := groups[group][r.GetInt(j, i, len(groups[group]))]
			tileRow[j] = Tile{texture: tile.Num}
		}
		tiles[i] = tileRow
	}
	return tiles
}

func genFill(tilesCanUse []*TileRender) [][]bool {
	return nil
}

func usableGroups(tilesCanUse []*TileRender) map[string][]*TileRender {
	usable := make(map[string][]*TileRender)
	for _, t := range tilesCanUse {
		if t.TileType == "M" {
			usable[t.Group] = append(usable[t.Group], t)
		}
	}
	return usable
}
