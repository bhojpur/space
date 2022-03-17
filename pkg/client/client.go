package client

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
	"github.com/bhojpur/space/pkg/client/base"
	"github.com/bhojpur/space/pkg/client/directions"
	directionsmatrix "github.com/bhojpur/space/pkg/client/directions_matrix"
	"github.com/bhojpur/space/pkg/client/geocode"
	mapmatching "github.com/bhojpur/space/pkg/client/map_matching"
	"github.com/bhojpur/space/pkg/client/maps"
)

// MapEngine API wrapper structure
type MapEngine struct {
	base *base.Base
	// Maps allows fetching of tiles and tilesets
	Maps *maps.Maps
	// Geocode allows forward (by address) and reverse (by lat/lng) geocoding
	Geocode *geocode.Geocode
	// Directions generates directions between arbitrary points
	Directions *directions.Directions
	// Direction Matrix returns all travel times and ways points between multiple points
	DirectionsMatrix *directionsmatrix.DirectionsMatrix
	// MapMatching snaps inaccurate path tracked to a map to produce a clean path
	MapMatching *mapmatching.MapMatching
}

// NewMapEngine Create a new MapEngine API instance
func NewMapEngine(token string) *MapEngine {
	m := &MapEngine{}

	// Create base instance
	m.base = base.NewBase(token)

	// Bind modules
	m.Maps = maps.NewMaps(m.base)
	m.Geocode = geocode.NewGeocode(m.base)
	m.Directions = directions.NewDirections(m.base)
	m.DirectionsMatrix = directionsmatrix.NewDirectionsMatrix(m.base)
	m.MapMatching = mapmatching.NewMapMaptching(m.base)

	return m
}
