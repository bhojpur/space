package geocode

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
	apiName    = "geocoding"
	apiVersion = "v1"
	apiMode    = "maps.places"
)

// Type defines geocode location response types
type Type string

const (
	// Country level
	Country Type = "country"
	// Region level
	Region Type = "region"
	// Postcode level
	Postcode Type = "postcode"
	// District level
	District Type = "district"
	// Place level
	Place Type = "place"
	// Locality level
	Locality Type = "locality"
	// Neighborhood level
	Neighborhood Type = "neighborhood"
	// Address level
	Address Type = "address"
	// POI (Point of Interest) level
	POI Type = "poi"
)

// Geocode api wrapper instance
type Geocode struct {
	base *base.Base
}

// NewGeocode Create a new Geocode API wrapper
func NewGeocode(base *base.Base) *Geocode {
	return &Geocode{base}
}

// ForwardRequestOpts request options fo forward geocoding
type ForwardRequestOpts struct {
	Country      string           `url:"country,omitempty"`
	Proximity    []float64        `url:"proximity,omitempty"`
	Types        []Type           `url:"types,omitempty"`
	Autocomplete bool             `url:"autocomplete,omitempty"`
	BBox         base.BoundingBox `url:"bbox,omitempty"`
	Limit        uint             `url:"limit,omitempty"`
}

// ForwardResponse is the response from a forward geocode lookup
type ForwardResponse struct {
	*base.FeatureCollection
	Query []string
}

// Forward geocode lookup
// Finds locations from a place name
func (g *Geocode) Forward(place string, req *ForwardRequestOpts) (*ForwardResponse, error) {

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ForwardResponse{}

	queryString := strings.Replace(place, " ", "+", -1)

	err = g.base.Query(apiName, apiVersion, apiMode, fmt.Sprintf("%s.json", queryString), &v, &resp)

	return &resp, err
}

// ReverseRequestOpts request options fo reverse geocoding
type ReverseRequestOpts struct {
	Types []Type
	Limit uint
}

// ReverseResponse is the response to a reverse geocode request
type ReverseResponse struct {
	*base.FeatureCollection
	Query []float64
}

// Reverse geocode lookup
// Finds place names from a location
func (g *Geocode) Reverse(loc *base.Location, req *ReverseRequestOpts) (*ReverseResponse, error) {

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ReverseResponse{}

	queryString := fmt.Sprintf("%f,%f.json", loc.Longitude, loc.Latitude)

	err = g.base.Query(apiName, apiVersion, apiMode, queryString, &v, &resp)

	return &resp, err
}
