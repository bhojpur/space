package geometry

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

// Geometry is a standard geometry
type Geometry interface {
	Rect() Rect
	Empty() bool
	Valid() bool
	ContainsPoint(point Point) bool
	IntersectsPoint(point Point) bool
	ContainsRect(rect Rect) bool
	IntersectsRect(rect Rect) bool
	ContainsLine(line *Line) bool
	IntersectsLine(line *Line) bool
	ContainsPoly(poly *Poly) bool
	IntersectsPoly(poly *Poly) bool
}

// require conformance
var _ = []Geometry{Point{}, Rect{}, &Line{}, &Poly{}}

// WorldPolygon is the maximum bounds for any GeoPoint
var WorldPolygon = NewPoly([]Point{
	{-180, -90}, {-180, 90}, {180, 90}, {180, -90}, {-180, -90},
}, nil, &IndexOptions{})
