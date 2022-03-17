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

// MapID selects which map to fetch from the API
type MapID string

// Map IDs
const (
	MapIDStreets          MapID = "maps.streets"
	MapIDLight            MapID = "maps.light"
	MapIDDark             MapID = "maps.dark"
	MapIDSatellite        MapID = "maps.satellite"
	MapIDStreetsSatellite MapID = "maps.streets-satellite"
	MapIDWheatpaste       MapID = "maps.wheatpaste"
	MapIDStreetsBasic     MapID = "maps.streets-basic"
	MapIDComic            MapID = "maps.comic"
	MapIDOutdoors         MapID = "maps.outdoors"
	MapIDRunBikeHike      MapID = "maps.run-bike-hike"
	MapIDPencil           MapID = "maps.pencil"
	MapIDPirates          MapID = "maps.pirates"
	MapIDEmerald          MapID = "maps.emerald"
	MapIDHighContrast     MapID = "maps.high-contrast"
	MapIDTerrainRGB       MapID = "maps.terrain-rgb"
)

// MapFormat specifies the format in which to return the map tiles
type MapFormat string

// Map formats
const (
	MapFormatPng        MapFormat = "png"    // true color PNG
	MapFormatPng32      MapFormat = "png32"  // 32 color indexed PNG
	MapFormatPng64      MapFormat = "png64"  // 64 color indexed PNG
	MapFormatPng128     MapFormat = "png128" // 128 color indexed PNG
	MapFormatPng256     MapFormat = "png256" // 256 color indexed PNG
	MapFormatPngRaw     MapFormat = "pngraw" // Raw PNG (only for MapIDTerrainRGB)
	MapFormatJpg70      MapFormat = "jpg70"  // 70% quality JPG
	MapFormatJpg80      MapFormat = "jpg80"  // 80% quality JPG
	MapFormatJpg90      MapFormat = "jpg90"  // 90% quality JPG
	MapFormatVectorTile MapFormat = "mvt"    // Vector Tile
)
