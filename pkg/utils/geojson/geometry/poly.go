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

// Poly ...
type Poly struct {
	Exterior Ring
	Holes    []Ring
}

// NewPoly ...
func NewPoly(exterior []Point, holes [][]Point, opts *IndexOptions) *Poly {
	poly := new(Poly)
	poly.Exterior = newRing(exterior, opts)
	if len(holes) > 0 {
		poly.Holes = make([]Ring, len(holes))
		for i := range holes {
			poly.Holes[i] = newRing(holes[i], opts)
		}
	}
	return poly
}

// Clockwise ...
func (poly *Poly) Clockwise() bool {
	if poly == nil || poly.Exterior == nil {
		return false
	}
	return poly.Exterior.Clockwise()
}

// Empty ...
func (poly *Poly) Empty() bool {
	if poly == nil || poly.Exterior == nil {
		return true
	}
	return poly.Exterior.Empty()
}

// Valid ...
func (poly *Poly) Valid() bool {
	if !poly.Exterior.Valid() {
		return false
	}
	for _, hole := range poly.Holes {
		if !hole.Valid() {
			return false
		}
	}
	return true
}

// Rect ...
func (poly *Poly) Rect() Rect {
	if poly == nil || poly.Exterior == nil {
		return Rect{}
	}
	return poly.Exterior.Rect()
}

// Move the polygon by delta. Returns a new polygon
func (poly *Poly) Move(deltaX, deltaY float64) *Poly {
	if poly == nil {
		return nil
	}
	if poly.Exterior == nil {
		return new(Poly)
	}
	npoly := new(Poly)
	if series, ok := poly.Exterior.(*baseSeries); ok {
		npoly.Exterior = Ring(series.Move(deltaX, deltaY))
	} else {
		nseries := makeSeries(
			seriesCopyPoints(poly.Exterior), false, true, DefaultIndexOptions)
		npoly.Exterior = Ring(nseries.Move(deltaX, deltaY))
	}
	if len(poly.Holes) > 0 {
		npoly.Holes = make([]Ring, len(poly.Holes))
		for i, hole := range poly.Holes {
			if series, ok := hole.(*baseSeries); ok {
				npoly.Holes[i] = Ring(series.Move(deltaX, deltaY))
			} else {
				nseries := makeSeries(
					seriesCopyPoints(hole), false, true, DefaultIndexOptions)
				npoly.Holes[i] = Ring(nseries.Move(deltaX, deltaY))
			}
		}
	}
	return npoly
}

// ContainsPoint ...
func (poly *Poly) ContainsPoint(point Point) bool {
	if poly == nil || poly.Exterior == nil {
		return false
	}
	if !ringContainsPoint(poly.Exterior, point, true).hit {
		return false
	}
	contains := true
	for _, hole := range poly.Holes {
		if ringContainsPoint(hole, point, false).hit {
			contains = false
			break
		}
	}
	return contains
}

// IntersectsPoint ...
func (poly *Poly) IntersectsPoint(point Point) bool {
	if poly == nil {
		return false
	}
	return poly.ContainsPoint(point)
}

// ContainsRect ...
func (poly *Poly) ContainsRect(rect Rect) bool {
	if poly == nil {
		return false
	}
	// convert rect into a polygon
	return poly.ContainsPoly(&Poly{Exterior: rect})
}

// IntersectsRect ...
func (poly *Poly) IntersectsRect(rect Rect) bool {
	if poly == nil {
		return false
	}
	// convert rect into a polygon
	return poly.IntersectsPoly(&Poly{Exterior: rect})
}

// ContainsLine ...
func (poly *Poly) ContainsLine(line *Line) bool {
	if poly == nil || poly.Exterior == nil || line == nil {
		return false
	}
	if !ringContainsLine(poly.Exterior, line, true) {
		return false
	}
	for _, polyHole := range poly.Holes {
		if ringIntersectsLine(polyHole, line, false) {
			return false
		}
	}
	return true
}

// IntersectsLine ...
func (poly *Poly) IntersectsLine(line *Line) bool {
	if poly == nil || poly.Exterior == nil || line == nil {
		return false
	}
	return ringIntersectsLine(poly.Exterior, line, true)
}

// ContainsPoly ...
func (poly *Poly) ContainsPoly(other *Poly) bool {
	if poly == nil || poly.Exterior == nil ||
		other == nil || other.Exterior == nil {
		return false
	}
	// 1) other exterior must be fully contained inside of the poly exterior.
	if !ringContainsRing(poly.Exterior, other.Exterior, true) {
		return false
	}
	// 2) ring cannot intersect poly holes
	contains := true
	for _, polyHole := range poly.Holes {
		if ringIntersectsRing(polyHole, other.Exterior, false) {
			contains = false
			// 3) unless the poly hole is contain inside of a other hole
			for _, otherHole := range other.Holes {
				if ringContainsRing(otherHole, polyHole, true) {
					contains = true
					// println(4)
					break
				}
			}
			if !contains {
				break
			}
		}
	}
	return contains
}

// IntersectsPoly ...
func (poly *Poly) IntersectsPoly(other *Poly) bool {
	if poly == nil || poly.Exterior == nil ||
		other == nil || other.Exterior == nil {
		return false
	}
	if !ringIntersectsRing(other.Exterior, poly.Exterior, true) {
		return false
	}
	for _, hole := range poly.Holes {
		if ringContainsRing(hole, other.Exterior, false) {
			return false
		}
	}
	for _, hole := range other.Holes {
		if ringContainsRing(hole, poly.Exterior, false) {
			return false
		}
	}
	return true
}
