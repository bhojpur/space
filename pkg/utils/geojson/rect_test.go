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
	"math/rand"
	"testing"

	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
)

func TestRect(t *testing.T) {
	rect := RO(10, 20, 30, 40)
	expect(t, !rect.Empty())
	expect(t, string(rect.AppendJSON(nil)) ==
		`{"type":"Polygon","coordinates":[[[10,20],[30,20],[30,40],[10,40],[10,20]]]}`)
	expect(t, rect.String() == string(rect.AppendJSON(nil)))
	// expect(t, !rect.Contains(NewString("")))
	// expect(t, !rect.Within(NewString("")))
	// expect(t, !rect.Intersects(NewString("")))
	// expect(t, rect.Distance(NewString("")) == 0)

	expect(t, rect.Rect() == R(10, 20, 30, 40))
	expect(t, rect.Center() == P(20, 30))
	var g Object
	rect.ForEach(func(o Object) bool {
		expect(t, g == nil)
		g = o
		return true
	})
	expect(t, g == rect)

	expect(t, rect.NumPoints() == 2)

	expect(t, !(&Point{}).Contains(rect))
	expect(t, !(&Rect{}).Contains(rect))
	expect(t, !(&LineString{}).Contains(rect))
	expect(t, !(&Polygon{}).Contains(rect))

	expect(t, !(&Point{}).Intersects(rect))
	expect(t, !(&Rect{}).Intersects(rect))
	expect(t, !(&LineString{}).Intersects(rect))
	expect(t, !(&Polygon{}).Intersects(rect))

	expect(t, (&Point{}).Distance(rect) != 0)
	expect(t, (&Rect{}).Distance(rect) != 0)
	expect(t, (&LineString{}).Distance(rect) != 0)
	expect(t, (&Polygon{}).Distance(rect) != 0)

}

func TestRectPoly(t *testing.T) {
	rect := RO(10, 20, 30, 40)
	json := rect.JSON()
	expect(t, json ==
		`{"type":"Polygon","coordinates":[[[10,20],[30,20],[30,40],[10,40],[10,20]]]}`)
	opts := *DefaultParseOptions
	opts.AllowRects = true
	o, err := Parse(json, &opts)
	expect(t, err == nil)
	rect2, ok := o.(*Rect)
	expect(t, ok)
	expect(t, rect2.base == rect.base)
	opts.AllowRects = false
	o, err = Parse(json, &opts)
	expect(t, err == nil)
	poly, ok := o.(*Polygon)
	expect(t, ok)
	json2 := poly.JSON()
	expect(t, json == json2)
}

func TestRectValid(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[[[10,200],[30,200],[30,40],[10,40],[10,200]]]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func BenchmarkRectValid(b *testing.B) {
	rects := make([]*Rect, b.N)
	for i := 0; i < b.N; i++ {
		min := geometry.Point{
			X: rand.Float64()*400 - 200, // some are out of bounds
			Y: rand.Float64()*200 - 100, // some are out of bounds
		}
		max := geometry.Point{
			X: rand.Float64()*400 - 200, // some are out of bounds
			Y: rand.Float64()*200 - 100, // some are out of bounds
		}
		if min.X > max.X {
			min.X, max.X = max.X, min.X
		}
		if min.Y > max.Y {
			min.Y, max.Y = max.Y, min.Y
		}
		rects[i] = NewRect(geometry.Rect{Min: min, Max: max})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rects[i].Valid()
	}
}
