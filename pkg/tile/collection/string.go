package collection

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
	"encoding/json"

	"github.com/bhojpur/space/pkg/utils/geojson"
	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
)

// String ...
type String string

var _ geojson.Object = String("")

// Spatial ...
func (s String) Spatial() geojson.Spatial {
	return geojson.EmptySpatial{}
}

// ForEach ...
func (s String) ForEach(iter func(geom geojson.Object) bool) bool {
	return iter(s)
}

// Empty ...
func (s String) Empty() bool {
	return true
}

// Valid ...
func (s String) Valid() bool {
	return false
}

// Rect ...
func (s String) Rect() geometry.Rect {
	return geometry.Rect{}
}

// Center ...
func (s String) Center() geometry.Point {
	return geometry.Point{}
}

// AppendJSON ...
func (s String) AppendJSON(dst []byte) []byte {
	data, _ := json.Marshal(string(s))
	return append(dst, data...)
}

// String ...
func (s String) String() string {
	return string(s)
}

// JSON ...
func (s String) JSON() string {
	return string(s.AppendJSON(nil))
}

// MarshalJSON ...
func (s String) MarshalJSON() ([]byte, error) {
	return s.AppendJSON(nil), nil
}

// Within ...
func (s String) Within(obj geojson.Object) bool {
	return false
}

// Contains ...
func (s String) Contains(obj geojson.Object) bool {
	return false
}

// Intersects ...
func (s String) Intersects(obj geojson.Object) bool {
	return false
}

// NumPoints ...
func (s String) NumPoints() int {
	return 0
}

// Distance ...
func (s String) Distance(obj geojson.Object) float64 {
	return 0
}
