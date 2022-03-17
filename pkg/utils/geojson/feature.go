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

import (
	"strings"

	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
	"github.com/bhojpur/space/pkg/utils/gjson"
	"github.com/bhojpur/space/pkg/utils/pretty"
	"github.com/bhojpur/space/pkg/utils/sjson"
)

// Feature ...
type Feature struct {
	base  Object
	extra *extra
}

// NewFeature returns a new GeoJSON Feature.
// The members must be a valid json object such as
// `{"id":"391","properties":{}}`, or it must be an empty string. It should not
// contain a "feature" member.
func NewFeature(geometry Object, members string) *Feature {
	g := new(Feature)
	g.base = geometry
	members = strings.TrimSpace(members)
	if members != "" && members != "{}" {
		if gjson.Valid(members) && gjson.Parse(members).IsObject() {
			if gjson.Get(members, "feature").Exists() {
				members, _ = sjson.Delete(members, "feature")
			}
			g.extra = new(extra)
			g.extra.members = string(pretty.UglyInPlace([]byte(members)))
		}
	}
	return g
}

// ForEach ...
func (g *Feature) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *Feature) Empty() bool {
	return g.base.Empty()
}

// Valid ...
func (g *Feature) Valid() bool {
	return g.base.Valid()
}

// Rect ...
func (g *Feature) Rect() geometry.Rect {
	return g.base.Rect()
}

// Center ...
func (g *Feature) Center() geometry.Point {
	return g.Rect().Center()
}

// Base ...
func (g *Feature) Base() Object {
	return g.base
}

// Members ...
func (g *Feature) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}

// AppendJSON ...
func (g *Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.base.AppendJSON(dst)
	dst = g.extra.appendJSONExtra(dst, true)
	dst = append(dst, '}')
	return dst

}

// String ...
func (g *Feature) String() string {
	return string(g.AppendJSON(nil))
}

// JSON ...
func (g *Feature) JSON() string {
	return string(g.AppendJSON(nil))
}

// MarshalJSON ...
func (g *Feature) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

// Spatial ...
func (g *Feature) Spatial() Spatial {
	return g
}

// Within ...
func (g *Feature) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Feature) Contains(obj Object) bool {
	return g.base.Contains(obj)
}

// WithinRect ...
func (g *Feature) WithinRect(rect geometry.Rect) bool {
	return g.base.Spatial().WithinRect(rect)
}

// WithinPoint ...
func (g *Feature) WithinPoint(point geometry.Point) bool {
	return g.base.Spatial().WithinPoint(point)
}

// WithinLine ...
func (g *Feature) WithinLine(line *geometry.Line) bool {
	return g.base.Spatial().WithinLine(line)
}

// WithinPoly ...
func (g *Feature) WithinPoly(poly *geometry.Poly) bool {
	return g.base.Spatial().WithinPoly(poly)
}

// Intersects ...
func (g *Feature) Intersects(obj Object) bool {
	return g.base.Intersects(obj)
}

// IntersectsPoint ...
func (g *Feature) IntersectsPoint(point geometry.Point) bool {
	return g.base.Spatial().IntersectsPoint(point)
}

// IntersectsRect ...
func (g *Feature) IntersectsRect(rect geometry.Rect) bool {
	return g.base.Spatial().IntersectsRect(rect)
}

// IntersectsLine ...
func (g *Feature) IntersectsLine(line *geometry.Line) bool {
	return g.base.Spatial().IntersectsLine(line)
}

// IntersectsPoly ...
func (g *Feature) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.Spatial().IntersectsPoly(poly)
}

// NumPoints ...
func (g *Feature) NumPoints() int {
	return g.base.NumPoints()
}

// parseJSONFeature will return a valid GeoJSON object.
func parseJSONFeature(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var g Feature
	if !keys.rGeometry.Exists() {
		return nil, errGeometryMissing
	}
	var err error
	g.base, err = Parse(keys.rGeometry.Raw, opts)
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	if point, ok := g.base.(*Point); ok {
		if g.extra != nil {
			members := g.extra.members
			if !opts.DisableCircleType &&
				gjson.Get(members, "properties.type").String() == "Circle" {
				// Circle
				radius := gjson.Get(members, "properties.radius").Float()
				units := gjson.Get(members, "properties.radius_units").String()
				switch units {
				case "", "m":
				case "km":
					radius *= 1000
				default:
					return nil, errCircleRadiusUnitsInvalid
				}
				return NewCircle(point.base, radius, 64), nil
			}
		}
	}
	return &g, nil
}

// Distance ...
func (g *Feature) Distance(obj Object) float64 {
	return g.base.Distance(obj)
}

// DistancePoint ...
func (g *Feature) DistancePoint(point geometry.Point) float64 {
	return g.base.Spatial().DistancePoint(point)
}

// DistanceRect ...
func (g *Feature) DistanceRect(rect geometry.Rect) float64 {
	return g.base.Spatial().DistanceRect(rect)
}

// DistanceLine ...
func (g *Feature) DistanceLine(line *geometry.Line) float64 {
	return g.base.Spatial().DistanceLine(line)
}

// DistancePoly ...
func (g *Feature) DistancePoly(poly *geometry.Poly) float64 {
	return g.base.Spatial().DistancePoly(poly)
}
