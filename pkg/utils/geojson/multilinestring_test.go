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

func TestMultiLineString(t *testing.T) {
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2,3]]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2]]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiLineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[1,null]}`, errCoordinatesInvalid)
}

func TestMultiLineStringValid(t *testing.T) {
	json := `{"type":"MultiLineString","coordinates":[
		[[10,10],[120,190]],
		[[50,50],[100,100]]
	]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestMultiLineStringPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiLineString","coordinates":[
		[[10,10],[20,20]],
		[[50,50],[100,100]]
	]}`, nil)
	expect(t, p.Intersects(PO(15, 15)))
	expect(t, p.Contains(PO(15, 15)))
	expect(t, p.Contains(PO(70, 70)))
	expect(t, !p.Contains(PO(40, 40)))
}
