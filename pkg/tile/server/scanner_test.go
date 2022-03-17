package server

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
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/bhojpur/space/pkg/utils/geojson"
	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
)

type testPointItem struct {
	object geojson.Object
	fields []float64
}

func PO(x, y float64) *geojson.Point {
	return geojson.NewPoint(geometry.Point{X: x, Y: y})
}

func BenchmarkFieldMatch(t *testing.B) {
	rand.Seed(time.Now().UnixNano())
	items := make([]testPointItem, t.N)
	for i := 0; i < t.N; i++ {
		items[i] = testPointItem{
			PO(rand.Float64()*360-180, rand.Float64()*180-90),
			[]float64{rand.Float64()*9 + 1, math.Round(rand.Float64()*30) + 1},
		}
	}
	sw := &scanWriter{
		wheres: []whereT{
			{"foo", 0, false, 1, false, 3},
			{"bar", 1, false, 10, false, 30},
		},
		whereins: []whereinT{
			{"foo", 0, []float64{1, 2}},
			{"bar", 1, []float64{11, 25}},
		},
		fmap: map[string]int{"foo": 0, "bar": 1},
		farr: []string{"bar", "foo"},
	}
	sw.fvals = make([]float64, len(sw.farr))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		// one call is super fast, measurements are not reliable, let's do 100
		for ix := 0; ix < 100; ix++ {
			sw.fieldMatch(items[i].fields, items[i].object)
		}
	}
}
