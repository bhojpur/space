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
	"os"
	"testing"

	"github.com/bhojpur/space/pkg/client/base"
	"github.com/bhojpur/space/pkg/client/directions"
	directionsmatrix "github.com/bhojpur/space/pkg/client/directions_matrix"
	mapmatching "github.com/bhojpur/space/pkg/client/map_matching"
	"github.com/bhojpur/space/pkg/client/geocode"
	"github.com/bhojpur/space/pkg/client/maps" // Import the core module and any required APIs
)

func TestMaps(t *testing.T) {
	// Fetch token from somewhere
	token := os.Getenv("BHOJPUR_SPACE_MAPS_TOKEN")
	if token == "" {
		t.Errorf("No token found")
		t.FailNow()
	}

	// Create new Bhojpur Space client instance
	engine := NewMapEngine(token)

	// Maps API
	_, err := engine.Maps.GetTile(maps.MapIDSatellite, 1, 0, 1, maps.MapFormatJpg90, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Geocoding API

	// Forward Geocoding
	var forwardOpts geocode.ForwardRequestOpts
	forwardOpts.Limit = 1

	place := "2 lincoln memorial circle nw"

	_, err = engine.Geocode.Forward(place, &forwardOpts)
	if err != nil {
		t.Error(err)
	}

	// Reverse Geocoding
	var reverseOpts geocode.ReverseRequestOpts
	reverseOpts.Limit = 1

	loc := &base.Location{72.438939, 34.074122}

	_, err = engine.Geocode.Reverse(loc, &reverseOpts)
	if err != nil {
		t.Error(err)
	}

	// Directions API
	var directionOpts directions.RequestOpts

	locs := []base.Location{{-122.42, 37.78}, {-77.03, 38.91}}

	_, err = engine.Directions.GetDirections(locs, directions.RoutingCycling, &directionOpts)
	if err != nil {
		t.Error(err)
	}

	// Directions Matrix API
	var directionMatrixOpts directionsmatrix.RequestOpts
	// Only 1st and second points will act as a source the response will be a 2x3 matrix
	source := []string{"0", "1"}
	dest := []string{"all"}
	directionMatrixOpts.SetSources(source)
	directionMatrixOpts.SetDestinations(dest)

	points := []base.Location{{37.752759, -122.467600}, {37.762819, -122.460304}, {37.758095, -122.442253}}

	_, err = engine.DirectionsMatrix.GetDirectionsMatrix(points, directionsmatrix.RoutingCycling, &directionMatrixOpts)
	if err != nil {
		t.Error(err)
	}

	//Map Matching API
	var MapMatchingOpts mapmatching.RequestOpts
	timeStamps := []int64{1492878132, 1492878142, 1492878152, 1492878172, 1492878182, 1492878192, 1492878202, 1492878302}
	radiusList := []int{9, 6, 8, 11, 8, 4, 8, 8}
	var opts mapmatching.RequestOpts
	opts.SetGeometries(mapmatching.GeometryPolyline)
	opts.SetOverview(mapmatching.OverviewFull)
	opts.SetTimestamps(timeStamps)
	opts.SetSteps(false)
	opts.SetAnnotations([]mapmatching.AnnotationType{mapmatching.AnnotationDistance, mapmatching.AnnotationSpeed})
	opts.SetRadiuses(radiusList)

	MatchingPath := []base.Location{{37.75319556403746, -122.44254112243651}, {37.75373846204306, -122.44238018989562},
		{37.754111702111146, -122.44199395179749}, {37.75473941979767, -122.44177401065825},
		{37.755570713402115, -122.4412429332733}, {37.756401997666046, -122.44113564491273},
		{37.75677098309616, -122.44228899478911}, {37.756949113334784, -122.4424821138382}}

	_, err = engine.MapMatching.GetMatching(MatchingPath, mapmatching.RoutingCycling, &MapMatchingOpts)
	if err != nil {
		t.Error(err)
	}
}
