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
	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
	"github.com/bhojpur/space/pkg/utils/gjson"
)

// MultiLineString ...
type MultiLineString struct{ collection }

// NewMultiLineString ...
func NewMultiLineString(lines []*geometry.Line) *MultiLineString {
	g := new(MultiLineString)
	for _, line := range lines {
		g.children = append(g.children, NewLineString(line))
	}
	g.parseInitRectIndex(DefaultParseOptions)
	return g
}

// AppendJSON ...
func (g *MultiLineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiLineString","coordinates":[`...)
	for i, g := range g.children {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = append(dst,
			gjson.GetBytes(g.AppendJSON(nil), "coordinates").String()...)
	}
	dst = append(dst, ']')
	if g.extra != nil {
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	return dst

}

// String ...
func (g *MultiLineString) String() string {
	return string(g.AppendJSON(nil))
}

// Valid ...
func (g *MultiLineString) Valid() bool {
	valid := true
	for _, p := range g.children {
		if !p.Valid() {
			valid = false
		}
	}
	return valid
}

// JSON ...
func (g *MultiLineString) JSON() string {
	return string(g.AppendJSON(nil))
}

// MarshalJSON ...
func (g *MultiLineString) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func parseJSONMultiLineString(
	keys *parseKeys, opts *ParseOptions,
) (Object, error) {
	var g MultiLineString
	var err error
	if !keys.rCoordinates.Exists() {
		return nil, errCoordinatesMissing
	}
	if !keys.rCoordinates.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords []geometry.Point
	var ex *extra
	keys.rCoordinates.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONLineStringCoords(keys, value, opts)
		if err != nil {
			return false
		}
		if len(coords) < 2 {
			err = errCoordinatesInvalid
			return false
		}
		gopts := toGeometryOpts(opts)
		line := geometry.NewLine(coords, &gopts)
		g.children = append(g.children, &LineString{base: *line, extra: ex})
		return true
	})
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	if opts.RequireValid {
		if !g.Valid() {
			return nil, errCoordinatesInvalid
		}
	}
	g.parseInitRectIndex(opts)
	return &g, nil
}
