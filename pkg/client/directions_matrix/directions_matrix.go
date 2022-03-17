package directionsmatrix

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
	"strings"

	"github.com/bhojpur/space/pkg/client/base"
	"github.com/google/go-querystring/query"
)

const (
	apiName    = "directions-matrix"
	apiVersion = "v1"
)

// DirectionsMatrix api wrapper instance
type DirectionsMatrix struct {
	base *base.Base
}

// NewDirectionsMatrix Create a new Directions Matrix API wrapper
func NewDirectionsMatrix(base *base.Base) *DirectionsMatrix {
	return &DirectionsMatrix{base}
}

// RoutingProfile defines routing mode for direction matrix finding
type RoutingProfile string

const (
	// RoutingDriving mode for for automovide routing
	RoutingDriving RoutingProfile = "maps/driving"
	// RoutingWalking mode for Pedestrian routing
	RoutingWalking RoutingProfile = "maps/walking"
	// RoutingCycling mode for bicycle routing
	RoutingCycling RoutingProfile = "maps/cycling"
)

// DirectionMatrixResponse is the response from GetDirections
type DirectionMatrixResponse struct {
	Code         string
	Durations    [][]float64
	Sources      []Waypoint
	Destinations []Waypoint
}

// Waypoint is an input point snapped to the road network
type Waypoint struct {
	Name     string
	Location []float64
}

// Codes are direction response Codes
type Codes string

const (
	// CodeOK success response
	CodeOK Codes = "Ok"
	//CodeProfileNotFound invalid routing profile
	CodeProfileNotFound Codes = "ProfileNotFound"
	// CodeInvalidInput invalid input data to the server
	CodeInvalidInput Codes = "InvalidInput"
)

// RequestOpts request options for directions api
type RequestOpts struct {
	Sources      string `url:"sources,omitempty"`
	Destinations string `url:"destinations,omitempty"`
}

// SetSources The points which will act as the starting point.
func (o *RequestOpts) SetSources(sources []string) {
	if sources[0] == "all" {
		o.Sources = "all"
	} else {
		lines := make([]string, len(sources))
		for i, r := range sources {
			lines[i] = fmt.Sprintf("%s", r)
		}
		o.Sources = strings.Join(lines, ";")
	}
}

// SetDestinations The points which will act as the destinations.
func (o *RequestOpts) SetDestinations(destinations []string) {
	if destinations[0] == "all" {
		o.Destinations = "all"
	} else {
		lines := make([]string, len(destinations))
		for i, r := range destinations {
			lines[i] = fmt.Sprintf("%s", r)
		}
		o.Destinations = strings.Join(lines, ";")
	}
}

// GetDirectionsMatrix between a set of locations using the specified routing profile
func (d *DirectionsMatrix) GetDirectionsMatrix(locations []base.Location, profile RoutingProfile, opts *RequestOpts) (*DirectionMatrixResponse, error) {

	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	coordinateStrings := make([]string, len(locations))
	for i, l := range locations {
		coordinateStrings[i] = fmt.Sprintf("%f,%f", l.Longitude, l.Latitude)
	}
	queryString := strings.Join(coordinateStrings, ";")

	resp := DirectionMatrixResponse{}

	err = d.base.Query(apiName, apiVersion, string(profile), queryString, &v, &resp)

	return &resp, err
}
