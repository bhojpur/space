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

// Point ...
type Point struct {
	X, Y float64
}

// Move ...
func (point Point) Move(deltaX, deltaY float64) Point {
	return Point{X: point.X + deltaX, Y: point.Y + deltaY}
}

// Empty ...
func (point Point) Empty() bool {
	return false
}

// Valid ...
func (point Point) Valid() bool {
	return point.X >= -180 && point.X <= 180 && point.Y >= -90 && point.Y <= 90
}

// Rect ...
func (point Point) Rect() Rect {
	return Rect{point, point}
}

// ContainsPoint ...
func (point Point) ContainsPoint(other Point) bool {
	return point == other
}

// IntersectsPoint ...
func (point Point) IntersectsPoint(other Point) bool {
	return point == other
}

// ContainsRect ...
func (point Point) ContainsRect(rect Rect) bool {
	return point.Rect() == rect
}

// IntersectsRect ...
func (point Point) IntersectsRect(rect Rect) bool {
	return rect.ContainsPoint(point)
}

// ContainsLine ...
func (point Point) ContainsLine(line *Line) bool {
	if line == nil {
		return false
	}
	return !line.Empty() && line.Rect() == point.Rect()
}

// IntersectsLine ...
func (point Point) IntersectsLine(line *Line) bool {
	if line == nil {
		return false
	}
	return line.IntersectsPoint(point)
}

// ContainsPoly ...
func (point Point) ContainsPoly(poly *Poly) bool {
	if poly == nil {
		return false
	}
	return !poly.Empty() && poly.Rect() == point.Rect()
}

// IntersectsPoly ...
func (point Point) IntersectsPoly(poly *Poly) bool {
	if poly == nil {
		return false
	}
	return poly.IntersectsPoint(point)
}
