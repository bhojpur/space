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

func TestMultiPoint(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3]]}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"MultiPoint","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2]],"bbox":null}`, nil)
	expectJSON(t, `{"type":"MultiPoint"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPoint","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3],[4,5,6]],"bbox":[1,2,3,4]}`, nil)
}

// func TestMultiPointPoly(t *testing.T) {
// 	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2],[2,2]]}`, nil)
// 	expect(t, p.Intersects(PO(1, 2)))
// 	expect(t, p.Contains(PO(1, 2)))
// 	expect(t, p.Contains(PO(2, 2)))
// 	expect(t, !p.Contains(PO(3, 2)))
// }
