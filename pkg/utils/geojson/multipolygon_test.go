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

import "testing"

func TestMultiPolygon(t *testing.T) {
	json := `{"type":"MultiPolygon","coordinates":[
		[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		],[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		]
	]}`
	p := expectJSON(t, json, nil)
	if p.Center() != P(5, 5) {
		t.Fatalf("expected '%v', got '%v'", P(5, 5), p.Center())
	}
	if cleanJSON(string(p.AppendJSON(nil))) != cleanJSON(json) {
		t.Fatalf("expectect '%v', got '%v'", cleanJSON(json), cleanJSON(string(p.AppendJSON(nil))))
	}
	expectJSON(t, `{"type":"MultiPolygon","coordinates":[[[[0,0],[10,0],[5,10],[0,0]],[[1,1]]]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPolygon"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPolygon","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPolygon","coordinates":[1,null]}`, errCoordinatesInvalid)
}

func TestMultiPolygonParseValid(t *testing.T) {
	json := `{"type":"MultiPolygon","coordinates":[
		[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		],[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		]
	]}`
	expectJSONOpts(t, json, nil, &ParseOptions{RequireValid: true})
}

func TestMultiPolygonPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiPolygon","coordinates":[
		[
			[[10,10],[20,20],[30,10],[10,10]]
		],[
			[[100,100],[200,200],[300,100],[100,100]]
		]
	]}`, nil)
	expect(t, p.Intersects(PO(15, 15)))
	expect(t, p.Contains(PO(15, 15)))
	expect(t, p.Contains(PO(150, 150)))
	expect(t, !p.Contains(PO(40, 40)))
	expect(t, p.Within(RO(-100, -100, 1000, 1000)))
}

func TestIssue369(t *testing.T) {
	poly14 := expectJSON(t, `{"type":"Polygon","coordinates":[[[-122.44154334068298,37.73179457567642],[-122.43935465812682,37.73179457567642],[-122.43935465812682,37.7343740514423],[-122.44154334068298,37.7343740514423],[-122.44154334068298,37.73179457567642]],[[-122.44104981422423,37.73286371140448],[-122.44104981422423,37.73424677678513],[-122.43990182876587,37.73424677678513],[-122.43990182876587,37.73286371140448],[-122.44104981422423,37.73286371140448]],[[-122.44109272956847,37.731870943026074],[-122.43976235389708,37.731870943026074],[-122.43976235389708,37.7326855231885],[-122.44109272956847,37.7326855231885],[-122.44109272956847,37.731870943026074]]]}`, nil)
	query := expectJSON(t, `{"type":"MultiPolygon","coordinates":[[[[-122.4408378,37.7341129],[-122.4408378,37.733],[-122.44,37.733],[-122.44,37.7341129],[-122.4408378,37.7341129]]],[[[-122.44091033935547,37.731981251280985],[-122.43994474411011,37.731981251280985],[-122.43994474411011,37.73254976045042],[-122.44091033935547,37.73254976045042],[-122.44091033935547,37.731981251280985]]]]}`, nil)
	expect(t, !query.Intersects(poly14))
	expect(t, !poly14.Intersects(query))
}
