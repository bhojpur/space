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

import "errors"

// LatLongToQuad iterates through all of the quads parts until levelOfDetail is reached.
func LatLongToQuad(latitude, longitude float64, levelOfDetail uint64, iterator func(part int) bool) {
	pixelX, pixelY := LatLongToPixelXY(latitude, longitude, levelOfDetail)
	tileX, tileY := PixelXYToTileXY(pixelX, pixelY)
	for i := levelOfDetail; i > 0; i-- {
		if !iterator(partForTileXY(tileX, tileY, i)) {
			break
		}
	}
}

func partForTileXY(tileX, tileY int64, levelOfDetail uint64) int {
	mask := int64(1 << (levelOfDetail - 1))
	if (tileX & mask) != 0 {
		if (tileY & mask) != 0 {
			return 3
		}
		return 1
	} else if (tileY & mask) != 0 {
		return 2
	}
	return 0
}

// TileXYToBounds returns the bounds around a tile.
func TileXYToBounds(tileX, tileY int64, levelOfDetail uint64) (minLat, minLon, maxLat, maxLon float64) {
	size := int64(1 << levelOfDetail)
	pixelX, pixelY := TileXYToPixelXY(tileX, tileY)
	maxLat, minLon = PixelXYToLatLong(pixelX, pixelY, levelOfDetail)
	pixelX, pixelY = TileXYToPixelXY(tileX+1, tileY+1)
	minLat, maxLon = PixelXYToLatLong(pixelX, pixelY, levelOfDetail)
	if size == 0 || tileX%size == 0 {
		minLon = MinLongitude
	}
	if size == 0 || tileX%size == size-1 {
		maxLon = MaxLongitude
	}
	if tileY <= 0 {
		maxLat = MaxLatitude
	}
	if tileY >= size-1 {
		minLat = MinLatitude
	}
	return
}

// QuadKeyToBounds converts a quadkey to bounds
func QuadKeyToBounds(quadkey string) (minLat, minLon, maxLat, maxLon float64, err error) {
	for i := 0; i < len(quadkey); i++ {
		switch quadkey[i] {
		case '0', '1', '2', '3':
		default:
			err = errors.New("invalid quadkey")
			return
		}
	}
	minLat, minLon, maxLat, maxLon = TileXYToBounds(QuadKeyToTileXY(quadkey))
	return
}
