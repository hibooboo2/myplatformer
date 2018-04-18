package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetInt(t *testing.T) {
	r := &rngSource{}
	for i := 1; i < 5; i++ {
		r.Seed(int(time.Now().UnixNano()))
		for x := -50; x < 50; x++ {
			for y := -50; y < 50; y++ {
				for max := 1; max < 20; max++ {
					a := r.GetInt(x, y, max)
					b := r.GetInt(x, y, max)
					if a != b {
						t.Logf("A:%v B:%v", a, b)
						t.Fail()
					}
				}
			}
		}
	}
}

func TestGetTexturesAndGenTiles(t *testing.T) {
	tiles, err := getTextures(nil)
	if err != nil {
		t.Log("Failed to get textures:", err)
		t.FailNow()
	}

	a := genTiles(20, 23425, tiles)
	b := genTiles(20, 23425, tiles)

	failed := 0
	for i := range a {
		for j := range a[i] {
			if a[i][j].texture != b[i][j].texture {
				failed++
				t.Fail()
			}
		}
	}
	if failed > 0 {
		t.Log("Tiles that failed to match:", failed)
		for i := range a {
			for j := range a[i] {
				fmt.Printf(" %v:%v ", a[i][j].texture, b[i][j].texture)
			}
			fmt.Printf("\n")
		}
	}
}
