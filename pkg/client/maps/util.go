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
	"bufio"
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/bhojpur/space/pkg/client/base"
)

// LocationToTileID converts a lat/lon location into a tile ID
func LocationToTileID(loc base.Location, level uint64) (uint64, uint64) {
	return MercatorLocationToTileID(loc.Latitude, loc.Longitude, level, 256)
}

// TileIDToLocation converts a tile ID to a lat/lon location
func TileIDToLocation(x, y float64, level uint64) base.Location {
	lat, lng := MercatorPixelToLocation(x, y, level, 256)
	return base.Location{
		Latitude:  lat,
		Longitude: lng,
	}
}

// WrapTileID wraps tile IDs by level for api requests
// eg. Tile (X:16, Y:10, level:4 ) will become (X:0, Y:10, level:4)
func WrapTileID(x, y, level uint64) (uint64, uint64) {
	// Limit to 2^n tile range for a given level
	x = x % (2 << (level - 1))
	y = y % (2 << (level - 1))

	return x, y
}

// GetEnclosingTileIDs fetches a pair of tile IDs enclosing the provided pair of points
func GetEnclosingTileIDs(a, b base.Location, level uint64) (uint64, uint64, uint64, uint64) {
	aX, aY := LocationToTileID(a, level)
	bX, bY := LocationToTileID(b, level)

	var xStart, xEnd, yStart, yEnd uint64
	if bX >= aX {
		xStart = aX
		xEnd = bX
	} else {
		xStart = bX
		xEnd = aX
	}

	if bY >= aY {
		yStart = aY
		yEnd = bY
	} else {
		yStart = bY
		yEnd = aY
	}

	return xStart, yStart, xEnd, yEnd
}

// StitchTiles combines a 2d array of image tiles into a single larger image
// Note that all images must have the same dimensions for this to work
func StitchTiles(images [][]Tile) Tile {

	imgX := images[0][0].Image.Bounds().Dx()
	imgY := images[0][0].Image.Bounds().Dy()

	xSize := imgX * len(images[0])
	ySize := imgY * len(images)

	stitched := image.NewRGBA(image.Rect(0, 0, xSize, ySize))

	for y, row := range images {
		for x, img := range row {
			sp := image.Point{0, 0}
			bounds := image.Rect(x*imgX, y*imgY, (x+1)*imgX, (y+1)*imgY)
			draw.Draw(stitched, bounds, img, sp, draw.Over)
		}
	}

	return NewTile(images[0][0].X, images[0][0].Y, images[0][0].Level, images[0][0].Size, stitched)
}

// LoadImage loads an image from a file
func LoadImage(file string) (image.Image, *image.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	r := bufio.NewReader(f)
	data, err := ioutil.ReadAll(r)
	f.Close()

	cfg := image.Config{}
	cfg, _, err = image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	return img, &cfg, nil
}

// SaveImageJPG writes an image instance to a jpg file
func SaveImageJPG(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

// SaveImagePNG writes an image instance to a png file
func SaveImagePNG(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	err = png.Encode(w, img)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

// PixelToHeight Converts a pixel to a height value for Bhojpur Space terrain tiles
func PixelToHeight(r, g, b uint8) float64 {
	R, G, B := float64(r), float64(g), float64(b)
	return -10000 + ((R*256*256 + G*256 + B) * 0.1)
}

func HeightToPixel(alt float64) (uint8, uint8, uint8) {
	increments := int((alt + 10000) / 0.1)
	b := uint8((increments >> 0) % 0xFF)
	g := uint8((increments >> 8) % 0xFF)
	r := uint8((increments >> 16) % 0xFF)
	return r, g, b
}
