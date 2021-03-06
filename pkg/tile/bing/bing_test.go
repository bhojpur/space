package bing

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
	"math/rand"
	"testing"
	"time"
)

func TestLevelFuzz(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		level := (rand.Int() % MaxLevelOfDetail) + 1
		quad := ""
		for j := 0; j < level; j++ {
			quad += string(byte(rand.Int()%4) + '0')
		}
		tileX, tileY, levelOfDetail := QuadKeyToTileXY(quad)
		if levelOfDetail != uint64(len(quad)) {
			t.Fatalf("[%d,%d] levelOfDetail == %d, expect %d", i, level, levelOfDetail, len(quad))
		}
		pixelX, pixelY := TileXYToPixelXY(tileX, tileY)
		latitude, longitude := PixelXYToLatLong(pixelX, pixelY, levelOfDetail)
		pixelX2, pixelY2 := LatLongToPixelXY(latitude, longitude, levelOfDetail)
		if pixelX2 != pixelX {
			t.Fatalf("[%d,%d] pixelX2 == %d, expect %d", i, level, pixelX2, pixelX)
		}
		if pixelY2 != pixelY {
			t.Fatalf("[%d,%d] pixelY2 == %d, expect %d", i, level, pixelY2, pixelY)
		}
		tileX2, tileY2 := PixelXYToTileXY(pixelX2, pixelY2)
		if tileX2 != tileX {
			t.Fatalf("[%d,%d] tileX2 == %d, expect %d", i, level, tileX2, tileX)
		}
		if tileY2 != tileY {
			t.Fatalf("[%d,%d] tileY2 == %d, expect %d", i, level, tileY2, tileY)
		}
		quad2 := TileXYToQuadKey(tileX2, tileY2, levelOfDetail)
		if quad2 != quad {
			t.Fatalf("[%d,%d] quad2 == %s, expect %s", i, level, quad2, quad)
		}
	}
}

func TestInvalidQuadKeyFuzz(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		func() {
			defer func() {
				var s string
				if v := recover(); v != nil {
					s = v.(string)
				}
				if s != "Invalid QuadKey digit sequence." {
					t.Fatalf("s == '%s', expect '%s", s, "Invalid QuadKey digit sequence.")
				}
			}()
			level := (rand.Int() % MaxLevelOfDetail) + 1

			valid := true
			quad := ""
			for valid {
				quad = ""
				for j := 0; j < level; j++ {
					c := byte(rand.Int()%5) + '0'
					quad += string(c)
					if c < '0' || c > '3' {
						valid = false
					}
				}
			}
			QuadKeyToTileXY(quad)
		}()
	}
}

func TestLatLonClippingFuzz(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		lat := clip(rand.Float64()*180.0-90.0, MinLatitude, MaxLatitude)
		lon := clip(rand.Float64()*380.0-180.0, MinLongitude, MaxLongitude)
		if lat < MinLatitude {
			t.Fatalf("lat == %f, expect < %f", lat, MinLatitude)
		}
		if lat > MaxLatitude {
			t.Fatalf("lat == %f, expect > %f", lat, MaxLatitude)
		}
		if lon < MinLongitude {
			t.Fatalf("lon == %f, expect < %f", lon, MinLongitude)
		}
		if lon > MaxLongitude {
			t.Fatalf("lon == %f, expect > %f", lon, MaxLongitude)
		}
	}
}

func TestIssue302(t *testing.T) {
	// Requesting tile with zoom level > 63 crashes the server #302
	for z := uint64(0); z < 256; z++ {
		tileX, tileY := PixelXYToTileXY(LatLongToPixelXY(33, -115, z))
		TileXYToBounds(tileX, tileY, z)
	}
}
