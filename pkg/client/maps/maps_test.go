package maps

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/space/pkg/client/base"
)

func TestMaps(t *testing.T) {

	token := os.Getenv("BHOJPUR_SERVICE_MAPS_TOKEN")
	if token == "" {
		t.Error("Maps API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)
	//b.SetDebug(true)

	maps := NewMaps(b)

	t.Run("Can fetch map tiles as png", func(t *testing.T) {

		img, err := maps.GetTile(MapIDStreets, 1, 0, 1, MapFormatPng, true)
		assert.Nil(t, err)

		err = SaveImagePNG(img, "/tmp/bhojpur-maps-test.png")
		assert.Nil(t, err)
	})

	t.Run("Can fetch map tiles as jpeg", func(t *testing.T) {

		img, err := maps.GetTile(MapIDSatellite, 1, 0, 1, MapFormatJpg90, true)
		assert.Nil(t, err)

		err = SaveImageJPG(img, "/tmp/bhojpur-maps-test.jpg")
		assert.Nil(t, err)
	})

	t.Run("Can fetch terrain RGB tiles", func(t *testing.T) {

		img, err := maps.GetTile(MapIDTerrainRGB, 1, 0, 1, MapFormatPngRaw, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		err = SaveImagePNG(img, "/tmp/bhojpur-maps-terrain.png")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("Can fetch map tiles by location", func(t *testing.T) {

		locA := base.Location{-45.942805, 166.568500}
		locB := base.Location{-34.2186101, 183.4015517}

		images, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 6, MapFormatJpg90, true)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		for y := range images {
			for x := range images[y] {
				SaveImageJPG(images[y][x], fmt.Sprintf("/tmp/bhojpur-maps-stitch-%d-%d.jpg", x, y))
			}
		}

		img := StitchTiles(images)

		err = SaveImageJPG(img, "/tmp/bhojpur-maps-stitch.jpg")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

	t.Run("Can fetch map tiles by location (with cache)", func(t *testing.T) {

		locA := base.Location{-45.942805, 166.568500}
		locB := base.Location{-34.2186101, 183.4015517}

		cache, err := NewFileCache("/tmp/bhojpur-maps-cache")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		maps.SetCache(cache)

		images, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 6, MapFormatJpg90, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		img := StitchTiles(images)

		err = SaveImageJPG(img, "/tmp/bhojpur-maps-stitch2.jpg")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

}
