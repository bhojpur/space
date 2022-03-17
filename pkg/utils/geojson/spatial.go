package geojson

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

import "github.com/bhojpur/space/pkg/utils/geojson/geometry"

// Spatial ...
type Spatial interface {
	WithinRect(rect geometry.Rect) bool
	WithinPoint(point geometry.Point) bool
	WithinLine(line *geometry.Line) bool
	WithinPoly(poly *geometry.Poly) bool
	IntersectsRect(rect geometry.Rect) bool
	IntersectsPoint(point geometry.Point) bool
	IntersectsLine(line *geometry.Line) bool
	IntersectsPoly(poly *geometry.Poly) bool
	DistanceRect(rect geometry.Rect) float64
	DistancePoint(point geometry.Point) float64
	DistanceLine(line *geometry.Line) float64
	DistancePoly(poly *geometry.Poly) float64
}

var _ = []Spatial{
	&Point{}, &LineString{}, &Polygon{}, &Feature{},
	&MultiPoint{}, &MultiLineString{}, &MultiPolygon{},
	&GeometryCollection{}, &FeatureCollection{}, &Rect{},
	EmptySpatial{},
}

// EmptySpatial ...
type EmptySpatial struct{}

// WithinRect ...
func (s EmptySpatial) WithinRect(rect geometry.Rect) bool {
	return false
}

// WithinPoint ...
func (s EmptySpatial) WithinPoint(point geometry.Point) bool {
	return false
}

// WithinLine ...
func (s EmptySpatial) WithinLine(line *geometry.Line) bool {
	return false
}

// WithinPoly ...
func (s EmptySpatial) WithinPoly(poly *geometry.Poly) bool {
	return false
}

// IntersectsRect ...
func (s EmptySpatial) IntersectsRect(rect geometry.Rect) bool {
	return false
}

// IntersectsPoint ...
func (s EmptySpatial) IntersectsPoint(point geometry.Point) bool {
	return false
}

// IntersectsLine ...
func (s EmptySpatial) IntersectsLine(line *geometry.Line) bool {
	return false
}

// IntersectsPoly ...
func (s EmptySpatial) IntersectsPoly(poly *geometry.Poly) bool {
	return false
}

// DistanceRect ...
func (s EmptySpatial) DistanceRect(rect geometry.Rect) float64 {
	return 0
}

// DistancePoint ...
func (s EmptySpatial) DistancePoint(point geometry.Point) float64 {
	return 0
}

// DistanceLine ...
func (s EmptySpatial) DistanceLine(line *geometry.Line) float64 {
	return 0
}

// DistancePoly ...
func (s EmptySpatial) DistancePoly(poly *geometry.Poly) float64 {
	return 0
}
