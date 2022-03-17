package mapmatching

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
	apiName    = "matching"
	apiVersion = "v5"
)

// MapMatching api wrapper instance
type MapMatching struct {
	base *base.Base
}

// NewMapMaptching Create a new Map Matching API wrapper
func NewMapMaptching(base *base.Base) *MapMatching {
	return &MapMatching{base}
}

// RoutingProfile defines routing mode for map matching
type RoutingProfile string

const (
	// RoutingDrivingTraffic mode for automotive routing takes into account current and historic traffic
	RoutingDrivingTraffic RoutingProfile = "maps/driving-traffic"
	// RoutingDriving mode for for automovide routing
	RoutingDriving RoutingProfile = "maps/driving"
	// RoutingWalking mode for Pedestrian routing
	RoutingWalking RoutingProfile = "maps/walking"
	// RoutingCycling mode for bicycle routing
	RoutingCycling RoutingProfile = "maps/cycling"
)

// RequestOpts request options for map matching api
type RequestOpts struct {
	Geometries  GeometryType    `url:"geometries,omitempty"`
	Radiuses    string          `url:"radiuses,omitempty"`
	Steps       bool            `url:"steps,omitempty"`
	Overview    OverviewType    `url:"overview,omitempty"`
	Timestamps  string          `url:"timestamps,omitempty"`
	Annotations *AnnotationType `url:"annotations,omitempty"`
}

// SetRadiuses sets radiuses for the maximum distance any coordinate can move when snapped to nearby road segment.
// This must have the same number of radiuses as locations in the GetMatching request
func (o *RequestOpts) SetRadiuses(radiuses []int) {
	lines := make([]string, len(radiuses))
	for i, r := range radiuses {
		lines[i] = fmt.Sprintf("%v", r)
	}
	o.Radiuses = strings.Join(lines, ";")
}

// SetAnnotations builds the annotations query argument from an array of annotation types
func (o *RequestOpts) SetAnnotations(annotations []AnnotationType) {
	lines := make([]string, len(annotations))
	for i, a := range annotations {
		lines[i] = fmt.Sprintf("%s", a)
	}
	o.Radiuses = strings.Join(lines, ",")
}

// SetTimestamps builds the Timestamps query argument from an array of timestamps types
// This must have the same number of timestamps as locations in the GetMatching request
func (o *RequestOpts) SetTimestamps(timestamps []int64) {
	lines := make([]string, len(timestamps))
	for i, a := range timestamps {
		lines[i] = fmt.Sprintf("%v", a)
	}
	o.Timestamps = strings.Join(lines, ";")
}

// SetGeometries builds the geometry query argument from the specified geometry type
func (o *RequestOpts) SetGeometries(geometrytype GeometryType) {
	o.Geometries = geometrytype
}

// SetOverview builds the overview query argument from the specified overview type
func (o *RequestOpts) SetOverview(overviewtype OverviewType) {
	o.Overview = overviewtype
}

// SetSteps builds the steps query argument from an array of steps option
func (o *RequestOpts) SetSteps(steps bool) {
	o.Steps = steps
}

// GetMatching for a path using the specified routing profile
func (d *MapMatching) GetMatching(path []base.Location, profile RoutingProfile, opts *RequestOpts) (*MatchingResponse, error) {

	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	coordinateStrings := make([]string, len(path))
	for i, l := range path {
		coordinateStrings[i] = fmt.Sprintf("%f,%f", l.Longitude, l.Latitude)
	}
	queryString := strings.Join(coordinateStrings, ";")

	resp := MatchingResponse{}

	err = d.base.Query(apiName, apiVersion, string(profile), queryString, &v, &resp)

	return &resp, err
}
