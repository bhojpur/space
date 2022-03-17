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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/space/pkg/client/base"
)

const delta = 1e-6

func TestMercator(t *testing.T) {
	zoom := uint64(4)
	size := uint64(256)
	fsize := float64(size)

	loc := base.Location{-45.942805, 166.568500}

	t.Run("Performs mercator projections to global pixels", func(t *testing.T) {
		x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)
		assert.EqualValues(t, 15.0, math.Floor(x/fsize))
		assert.EqualValues(t, 10.0, math.Floor(y/fsize))

		// Increase zoom scale x2 multiplies location by 4
		x, y = MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom+2, size)
		assert.EqualValues(t, 15.0*4+1, math.Floor(x/fsize))
		assert.EqualValues(t, 10.0*4+1, math.Floor(y/fsize))

		// Doubling tile size doubles pixel location
		x, y = MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size*2)
		assert.EqualValues(t, 15.0*2, math.Floor(x/fsize))
		assert.EqualValues(t, 10.0*2, math.Floor(y/fsize))
	})

	t.Run("Performs mercator projections to tile IDs", func(t *testing.T) {
		x, y := MercatorLocationToTileID(loc.Latitude, loc.Longitude, zoom, size)
		assert.EqualValues(t, 15, x)
		assert.EqualValues(t, 10, y)

		// Increasing zoom level by 2 multiplies tile IDs by 4
		x, y = MercatorLocationToTileID(loc.Latitude, loc.Longitude, zoom+2, size)
		assert.EqualValues(t, 15*4+1, x)
		assert.EqualValues(t, 10*4+1, y)

		// Doubling tile size does not change tile ID
		x, y = MercatorLocationToTileID(loc.Latitude, loc.Longitude, zoom, size*2)
		assert.EqualValues(t, 15, x)
		assert.EqualValues(t, 10, y)
	})

	t.Run("Reverses mercator projections", func(t *testing.T) {
		x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)

		lat2, lng2 := MercatorPixelToLocation(x, y, zoom, size)
		assert.InDelta(t, loc.Latitude, lat2, delta)
		assert.InDelta(t, loc.Longitude, lng2, delta)
	})
}

func BenchmarkMercator(b *testing.B) {
	zoom := uint64(4)
	loc := base.Location{Latitude: -45.942805, Longitude: 166.568500}
	size := uint64(256)

	x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)

	b.Run("Forward projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)
		}
	})

	b.Run("Reverse projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MercatorPixelToLocation(x, y, zoom, size)
		}
	})

}
