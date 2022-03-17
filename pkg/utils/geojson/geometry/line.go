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

// Line is a open series of points
type Line struct {
	baseSeries
}

// Valid ...
func (line *Line) Valid() bool {
	return line.baseSeries.Valid()
}

// NewLine creates a new Line
func NewLine(points []Point, opts *IndexOptions) *Line {
	line := new(Line)
	line.baseSeries = makeSeries(points, true, false, opts)
	return line
}

// Move ...
func (line *Line) Move(deltaX, deltaY float64) *Line {
	if line == nil {
		return nil
	}
	nline := new(Line)
	nline.baseSeries = *line.baseSeries.Move(deltaX, deltaY).(*baseSeries)
	return nline
}

// ContainsPoint ...
func (line *Line) ContainsPoint(point Point) bool {
	if line == nil {
		return false
	}
	contains := false
	line.Search(Rect{point, point}, func(seg Segment, index int) bool {
		if seg.Raycast(point).On {
			contains = true
			return false
		}
		return true
	})
	return contains
}

// IntersectsPoint ...
func (line *Line) IntersectsPoint(point Point) bool {
	if line == nil {
		return false
	}
	return line.ContainsPoint(point)
}

// ContainsRect ...
func (line *Line) ContainsRect(rect Rect) bool {
	if line == nil {
		return false
	}
	// Convert rect into a poly
	return line.ContainsPoly(&Poly{Exterior: rect})
}

// IntersectsRect ...
func (line *Line) IntersectsRect(rect Rect) bool {
	if line == nil {
		return false
	}
	return rect.IntersectsLine(line)
}

// ContainsLine ...
func (line *Line) ContainsLine(other *Line) bool {
	if line == nil || other == nil || line.Empty() || other.Empty() {
		return false
	}
	// locate the first "other" segment that contains the first "line" segment.
	lineNumSegments := line.NumSegments()
	segIdx := -1
	for j := 0; j < lineNumSegments; j++ {
		if line.SegmentAt(j).ContainsSegment(other.SegmentAt(0)) {
			segIdx = j
			break
		}
	}
	if segIdx == -1 {
		return false
	}
	otherNumSegments := other.NumSegments()
	for i := 1; i < otherNumSegments; i++ {
		lineSeg := line.SegmentAt(segIdx)
		otherSeg := other.SegmentAt(i)
		if lineSeg.ContainsSegment(otherSeg) {
			continue
		}
		if otherSeg.A == lineSeg.A {
			// reverse it
			if segIdx == 0 {
				return false
			}
			segIdx--
			i--
		} else if otherSeg.A == lineSeg.B {
			// forward it
			if segIdx == lineNumSegments-1 {
				return false
			}
			segIdx++
			i--
		}
	}
	return true
}

// IntersectsLine ...
func (line *Line) IntersectsLine(other *Line) bool {
	if line == nil || other == nil || line.Empty() || other.Empty() {
		return false
	}
	if !line.Rect().IntersectsRect(other.Rect()) {
		return false
	}
	if line.NumPoints() > other.NumPoints() {
		line, other = other, line
	}
	lineNumSegments := line.NumSegments()
	for i := 0; i < lineNumSegments; i++ {
		segA := line.SegmentAt(i)
		var intersects bool
		other.Search(segA.Rect(), func(segB Segment, _ int) bool {
			if segA.IntersectsSegment(segB) {
				intersects = true
				return false
			}
			return true
		})
		if intersects {
			return true
		}
	}
	return false
}

// ContainsPoly ...
func (line *Line) ContainsPoly(poly *Poly) bool {
	if line == nil || poly == nil || line.Empty() || poly.Empty() {
		return false
	}
	rect := poly.Rect()
	if rect.Min.X != rect.Max.X && rect.Min.Y != rect.Max.Y {
		return false
	}
	// polygon can fit in a straight (vertial or horizontal) line
	points := [2]Point{rect.Min, rect.Max}
	var other Line
	other.baseSeries.points = points[:]
	other.baseSeries.rect = rect
	return line.ContainsLine(&other)
}

// IntersectsPoly ...
func (line *Line) IntersectsPoly(poly *Poly) bool {
	return poly.IntersectsLine(line)
}
