package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Tile struct {
	Name    string
	texture int
	Loc     sdl.Point
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
