package clip

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
	"github.com/bhojpur/space/pkg/utils/geojson"
	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
)

func clipLineString(
	lineString *geojson.LineString, clipper geojson.Object,
	opts *geometry.IndexOptions,
) geojson.Object {
	bbox := clipper.Rect()
	var newPoints [][]geometry.Point
	var clipped geometry.Segment
	var rejected bool
	var line []geometry.Point
	base := lineString.Base()
	nSegments := base.NumSegments()
	for i := 0; i < nSegments; i++ {
		clipped, rejected = clipSegment(base.SegmentAt(i), bbox)
		if rejected {
			continue
		}
		if len(line) > 0 && line[len(line)-1] != clipped.A {
			newPoints = append(newPoints, line)
			line = []geometry.Point{clipped.A}
		} else if len(line) == 0 {
			line = append(line, clipped.A)
		}
		line = append(line, clipped.B)
	}
	if len(line) > 0 {
		newPoints = append(newPoints, line)
	}
	var children []*geometry.Line
	for _, points := range newPoints {
		children = append(children,
			geometry.NewLine(points, opts))
	}
	if len(children) == 1 {
		return geojson.NewLineString(children[0])
	}
	return geojson.NewMultiLineString(children)
}
