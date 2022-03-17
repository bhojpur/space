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

import "math"

// RaycastResult holds the results of the Raycast operation
type RaycastResult struct {
	In bool // point on the left
	On bool // point is directly on top of
}

// Raycast performs the raycast operation
func (seg Segment) Raycast(point Point) RaycastResult {

	p, a, b := point, seg.A, seg.B
	// make sure that the point is inside the segment bounds
	if a.Y < b.Y && (p.Y < a.Y || p.Y > b.Y) {
		return RaycastResult{false, false}
	} else if a.Y > b.Y && (p.Y < b.Y || p.Y > a.Y) {
		return RaycastResult{false, false}
	}

	// test if point is in on the segment
	if a.Y == b.Y {
		if a.X == b.X {
			if p == a {
				return RaycastResult{false, true}
			}
			return RaycastResult{false, false}
		}
		if p.Y == b.Y {
			// horizontal segment
			// check if the point in on the line
			if a.X < b.X {
				if p.X >= a.X && p.X <= b.X {
					return RaycastResult{false, true}
				}
			} else {
				if p.X >= b.X && p.X <= a.X {
					return RaycastResult{false, true}
				}
			}
		}
	}
	if a.X == b.X && p.X == b.X {
		// vertical segment
		// check if the point in on the line
		if a.Y < b.Y {
			if p.Y >= a.Y && p.Y <= b.Y {
				return RaycastResult{false, true}
			}
		} else {
			if p.Y >= b.Y && p.Y <= a.Y {
				return RaycastResult{false, true}
			}
		}
	}
	if (p.X-a.X)/(b.X-a.X) == (p.Y-a.Y)/(b.Y-a.Y) {
		return RaycastResult{false, true}
	}

	// do the actual raycast here.
	for p.Y == a.Y || p.Y == b.Y {
		p.Y = math.Nextafter(p.Y, math.Inf(1))
	}
	if a.Y < b.Y {
		if p.Y < a.Y || p.Y > b.Y {
			return RaycastResult{false, false}
		}
	} else {
		if p.Y < b.Y || p.Y > a.Y {
			return RaycastResult{false, false}
		}
	}
	if a.X > b.X {
		if p.X >= a.X {
			return RaycastResult{false, false}
		}
		if p.X <= b.X {
			return RaycastResult{true, false}
		}
	} else {
		if p.X >= b.X {
			return RaycastResult{false, false}
		}
		if p.X <= a.X {
			return RaycastResult{true, false}
		}
	}
	if a.Y < b.Y {
		if (p.Y-a.Y)/(p.X-a.X) >= (b.Y-a.Y)/(b.X-a.X) {
			return RaycastResult{true, false}
		}
	} else {
		if (p.Y-b.Y)/(p.X-b.X) >= (a.Y-b.Y)/(a.X-b.X) {
			return RaycastResult{true, false}
		}
	}
	return RaycastResult{false, false}
}
