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
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
	"github.com/bhojpur/space/pkg/utils/lotsa"
)

func parseCollection(t *testing.T, data string, index bool) Collection {
	t.Helper()
	opts := *DefaultParseOptions
	if index {
		opts.IndexChildren = 1
		opts.IndexGeometry = 1
	} else {
		opts.IndexChildren = 0
		opts.IndexGeometry = 0
	}
	g, err := Parse(string(data), &opts)
	if err != nil {
		t.Fatal(err)
	}
	c, ok := g.(Collection)
	if !ok {
		t.Fatal("not searchable")
	}
	return c
}

func testCollectionSubset(
	t *testing.T, json string, subsetRect geometry.Rect,
) []Collection {
	t.Helper()
	expectJSON(t, json, nil)
	start := time.Now()
	c := parseCollection(t, json, true)
	dur := time.Since(start)
	fmt.Printf("%d children in %s\n", len(c.Children()), dur)
	cs := make([]Collection, 2)
	ts := make([]time.Duration, 2)
	start = time.Now()
	cs[0] = parseCollection(t, json, false)
	ts[0] = time.Since(start)
	start = time.Now()
	cs[1] = parseCollection(t, json, true)
	ts[1] = time.Since(start)
	var lastSubset string
	for i, c := range cs {
		var children []Object
		start := time.Now()
		c.Search(subsetRect, func(child Object) bool {
			children = append(children, child)
			return true
		})
		dur := time.Since(start)
		if c.Indexed() {
			fmt.Printf("Indexed: %v (build: %v)\n", dur, ts[i])
		} else {
			fmt.Printf("Simple:  %v (build: %v)\n", dur, ts[i])
		}
		var childrenJSONs []string
		for _, child := range children {
			childrenJSONs = append(childrenJSONs, string(child.AppendJSON(nil)))
		}
		sort.Strings(childrenJSONs)
		subset := `{"type":"GeometryCollection","geometries":[` +
			strings.Join(childrenJSONs, ",") + `]}`
		if i > 0 {
			if subset != lastSubset {
				t.Fatal("mismatch")
			}
		}
		lastSubset = subset
	}
	return cs
}

func TestCollectionBostonSubset(t *testing.T) {
	data, err := ioutil.ReadFile("test_files/boston_subset.geojson")
	expect(t, err == nil)
	cs := testCollectionSubset(t, string(data),
		R(-71.474046, 42.492479, -71.466321, 42.497415),
	)
	expect(t, cs[0].(Object).Intersects(PO(-71.46723, 42.49432)))
	expect(t, cs[1].(Object).Intersects(PO(-71.46723, 42.49432)))
	expect(t, !cs[0].(*FeatureCollection).Intersects(PO(-71.46713, 42.49431)))
	expect(t, !cs[1].(Object).Intersects(PO(-71.46713, 42.49431)))
}

func TestCollectionUSASubset(t *testing.T) {
	data, err := ioutil.ReadFile("test_files/usa.geojson")
	if err != nil {
		fmt.Printf("test_files/usa.geojson missing\n")
		return
	}

	expect(t, err == nil)
	cs := testCollectionSubset(t, string(data),
		R(-90, -45, 90, 45),
	)
	expect(t, cs[0].(Object).Intersects(PO(-91.09863, 30.03105)))
	expect(t, cs[1].(Object).Intersects(PO(-91.09863, 30.03105)))
	expect(t, !cs[0].(Object).Intersects(PO(-86.87988, 28.439713)))
	expect(t, !cs[1].(Object).Intersects(PO(-86.87988, 28.439713)))

	N := 10000
	T := 6

	// 34 is the lower 48
	rect := cs[0].Children()[34].(Object).Rect()
	points := make([]Object, N)
	for i := 0; i < N; i++ {
		points[i] = PO(
			(rect.Max.X-rect.Min.X)*rand.Float64()+rect.Min.X,
			(rect.Max.Y-rect.Min.Y)*rand.Float64()+rect.Min.Y,
		)
	}
	lotsa.Output = os.Stdout
	cs0 := cs[0].(Object)
	cs1 := cs[1].(Object)
	lotsa.Ops(N, T, func(i, _ int) {
		cs0.Intersects(points[i])
	})
	lotsa.Ops(N, T, func(i, _ int) {
		cs1.Intersects(points[i])
	})

}
