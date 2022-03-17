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
	"testing"
)

func TestFeatureCollection(t *testing.T) {
	p := expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2,3]}]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"FeatureCollection"}`, errFeaturesMissing)
	expectJSON(t, `{"type":"FeatureCollection","features":null}`, errFeaturesInvalid)
	expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2,3]}],"bbox":null}`, nil)
	expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point"}]}`, errCoordinatesMissing)
}

func TestFeatureCollectionPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2]}]}`, nil)
	expect(t, p.Intersects(PO(1, 2)))
	expect(t, p.Contains(PO(1, 2)))
}

func TestFeatureCollectionValid(t *testing.T) {
	json := `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,200]}]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestForEach(t *testing.T) {
	json := `{"type":"FeatureCollection","features":[
		{"type":"Feature","id":"A","geometry":{"type":"Point","coordinates":[1,2]},"properties":{}},
		{"type":"Feature","id":"B","geometry":{"type":"Point","coordinates":[3,4]},"properties":{}},
		{"type":"Feature","id":"C","geometry":{"type":"Point","coordinates":[5,6]},"properties":{}},
		{"type":"Feature","id":"D","geometry":{"type":"Point","coordinates":[7,8]},"properties":{}}
	]}`

	g, _ := Parse(json, nil)
	objsA := g.(*FeatureCollection).Children()
	var objsB []Object
	g.ForEach(func(geom Object) bool {
		objsB = append(objsB, geom)
		return true
	})
	for i := 0; i < len(objsA) && i < len(objsB); i++ {
		expect(t, objsA[i].String() == objsB[i].String())
	}
}
