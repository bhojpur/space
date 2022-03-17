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

// Clip clips the contents of a geojson object and return
func Clip(
	obj geojson.Object, clipper geojson.Object, opts *geometry.IndexOptions,
) (clipped geojson.Object) {
	switch obj := obj.(type) {
	case *geojson.Point:
		return clipPoint(obj, clipper, opts)
	case *geojson.Rect:
		return clipRect(obj, clipper, opts)
	case *geojson.LineString:
		return clipLineString(obj, clipper, opts)
	case *geojson.Polygon:
		return clipPolygon(obj, clipper, opts)
	case *geojson.Feature:
		return clipFeature(obj, clipper, opts)
	case geojson.Collection:
		return clipCollection(obj, clipper, opts)
	}
	return obj
}

// clipSegment is Cohen-Sutherland Line Clipping
func clipSegment(seg geometry.Segment, rect geometry.Rect) (
	res geometry.Segment, rejected bool,
) {
	startCode := getCode(rect, seg.A)
	endCode := getCode(rect, seg.B)
	if (startCode | endCode) == 0 {
		// trivially accept
		res = seg
	} else if (startCode & endCode) != 0 {
		// trivially reject
		rejected = true
	} else if startCode != 0 {
		// start is outside. get new start.
		newStart := intersect(rect, startCode, seg.A, seg.B)
		res, rejected =
			clipSegment(geometry.Segment{A: newStart, B: seg.B}, rect)
	} else {
		// end is outside. get new end.
		newEnd := intersect(rect, endCode, seg.A, seg.B)
		res, rejected = clipSegment(geometry.Segment{A: seg.A, B: newEnd}, rect)
	}
	return
}

// clipRing is Sutherland-Hodgman Polygon Clipping
// https://www.cs.helsinki.fi/group/goa/viewing/leikkaus/intro2.html
func clipRing(ring []geometry.Point, bbox geometry.Rect) (
	resRing []geometry.Point,
) {
	if len(ring) < 4 {
		// under 4 elements this is not a polygon ring!
		return
	}
	var edge uint8
	var inside, prevInside bool
	var prev geometry.Point
	for edge = 1; edge <= 8; edge *= 2 {
		prev = ring[len(ring)-2]
		prevInside = (getCode(bbox, prev) & edge) == 0
		for _, p := range ring {
			inside = (getCode(bbox, p) & edge) == 0
			if prevInside && inside {
				// Staying inside
				resRing = append(resRing, p)
			} else if prevInside && !inside {
				// Leaving
				resRing = append(resRing, intersect(bbox, edge, prev, p))
			} else if !prevInside && inside {
				// Entering
				resRing = append(resRing, intersect(bbox, edge, prev, p))
				resRing = append(resRing, p)
			} /* else {
				// Stay outside
			} */
			prev, prevInside = p, inside
		}
		if len(resRing) > 0 && resRing[0] != resRing[len(resRing)-1] {
			resRing = append(resRing, resRing[0])
		}
		ring, resRing = resRing, []geometry.Point{}
		if len(ring) == 0 {
			break
		}
	}
	resRing = ring
	return
}

func getCode(bbox geometry.Rect, point geometry.Point) (code uint8) {
	code = 0

	if point.X < bbox.Min.X {
		code |= 1 // left
	} else if point.X > bbox.Max.X {
		code |= 2 // right
	}

	if point.Y < bbox.Min.Y {
		code |= 4 // bottom
	} else if point.Y > bbox.Max.Y {
		code |= 8 // top
	}

	return
}

func intersect(bbox geometry.Rect, code uint8, start, end geometry.Point) (
	new geometry.Point,
) {
	if (code & 8) != 0 { // top
		new = geometry.Point{
			X: start.X + (end.X-start.X)*(bbox.Max.Y-start.Y)/(end.Y-start.Y),
			Y: bbox.Max.Y,
		}
	} else if (code & 4) != 0 { // bottom
		new = geometry.Point{
			X: start.X + (end.X-start.X)*(bbox.Min.Y-start.Y)/(end.Y-start.Y),
			Y: bbox.Min.Y,
		}
	} else if (code & 2) != 0 { //right
		new = geometry.Point{
			X: bbox.Max.X,
			Y: start.Y + (end.Y-start.Y)*(bbox.Max.X-start.X)/(end.X-start.X),
		}
	} else if (code & 1) != 0 { // left
		new = geometry.Point{
			X: bbox.Min.X,
			Y: start.Y + (end.Y-start.Y)*(bbox.Min.X-start.X)/(end.X-start.X),
		}
	} /* else {
		// should not call intersect with the zero code
	} */

	return
}
