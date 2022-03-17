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
	"math"
)

const (
	// D2R helper for converting degrees to radians
	D2R = math.Pi / 180
	// R2D helper for converting radians to degrees
	R2D = 180 / math.Pi

	pi = math.Pi
)

// MercatorLocationToPixel converts a lat/lng/zoom location (in degrees) to a pixel location in the global space
func MercatorLocationToPixel(lat, lng float64, zoom, size uint64) (float64, float64) {
	latRad, lngRad := lat*D2R, lng*D2R
	fsize := float64(size)
	x := (fsize / 2 / pi) * math.Pow(2, float64(zoom)) * (lngRad + pi)
	y := (fsize / 2 / pi) * math.Pow(2, float64(zoom)) * (pi - math.Log(math.Tan(pi/4+latRad/2)))
	return x, y
}

// MercatorLocationToTileID builds on MercatorLocationToPixel to fetch the TileID of a given location
func MercatorLocationToTileID(lat, lng float64, zoom, size uint64) (uint64, uint64) {
	fsize := float64(size)
	x, y := MercatorLocationToPixel(lat, lng, zoom, size)
	xID, yID := uint64(x/fsize), uint64(y/fsize)
	return xID, yID
}

// MercatorPixelToLocation converts a given (global) pixel location and zoom level to a lat and lng (in degrees)
func MercatorPixelToLocation(x, y float64, zoom, size uint64) (float64, float64) {
	fsize := float64(size)
	lng := x*(pi/fsize*2)/math.Pow(2, float64(zoom)) - pi
	lat := (math.Atan(math.Pow(math.E, (-y*(pi/fsize*2)/math.Pow(2, float64(zoom))+pi))) - pi/4) * 2
	return lat * R2D, lng * R2D
}
