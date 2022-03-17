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

func TestLineStringParse(t *testing.T) {
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2,3]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"LineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2,3,4,5]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1]]}`, errCoordinatesInvalid)
	g := expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]]}`, nil)
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]],"bbox":null}`, nil)
	expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]],"bbox":[1,2,3,4]}`, nil)
}

func TestLineStringParseValid(t *testing.T) {
	json := `{"type":"LineString","coordinates":[[1,2],[-12,-190]]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errDataInvalid, &ParseOptions{RequireValid: true})
}

func TestLineStringVarious(t *testing.T) {
	var g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]]}`, nil)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"LineString","coordinates":[[3,4],[1,2]]}`)
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expect(t, g.Center() == P(2, 3))
	expect(t, !g.Empty())
	g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]],"bbox":[1,2,3,4]}`, nil)
	expect(t, !g.Empty())
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expect(t, g.Center() == R(1, 2, 3, 4).Center())
}

func TestLineStringValid(t *testing.T) {
	var g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]]}`, nil)
	expect(t, g.Valid())
}

func TestLineStringInvalid(t *testing.T) {
	var g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2],[0, 190]]}`, nil)
	expect(t, !g.Valid())
}

// func TestLineStringPoly(t *testing.T) {
// 	ls := expectJSON(t, `{"type":"LineString","coordinates":[
// 		[10,10],[20,20],[20,10]
// 	]}`, nil)
// 	expect(t, ls.(*LineString).Contains(ls))
// 	expect(t, ls.Contains(PO(10, 10)))
// 	expect(t, ls.Contains(PO(15, 15)))
// 	expect(t, ls.Contains(PO(20, 20)))
// 	expect(t, ls.Contains(PO(20, 15)))
// 	expect(t, !ls.Contains(PO(12, 13)))
// 	expect(t, !ls.Contains(RO(10, 10, 20, 20)))
// 	expect(t, ls.Intersects(PO(10, 10)))
// 	expect(t, ls.Intersects(PO(15, 15)))
// 	expect(t, ls.Intersects(PO(20, 20)))
// 	expect(t, !ls.Intersects(PO(12, 13)))
// 	expect(t, ls.Intersects(RO(10, 10, 20, 20)))
// 	expect(t, ls.Intersects(
// 		expectJSON(t, `{"type":"Point","coordinates":[15,15,0]}`, nil),
// 	))
// 	expect(t, ls.Intersects(ls))
// 	lsb := expectJSON(t, `{"type":"LineString","coordinates":[
// 		[10,10],[20,20],[20,10]
// 	],"bbox":[10,10,20,20]}`, nil)
// 	expect(t, lsb.Contains(PO(12, 13)))
// 	expect(t, ls.Contains(PO(20, 20)))
// }
