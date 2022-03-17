package algo

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

// Box performs simple box-distance algorithm on rectangles. When wrapX
// is provided, the operation does a cylinder wrapping of the X value to allow
// for antimeridian calculations. When itemDist is provided (not nil), it
// becomes the caller's responsibility to return the box-distance.
func Box(
	targetMin, targetMax [2]float64, wrapX bool,
	itemDist func(min, max [2]float64, data interface{}) float64,
) (
	algo func(min, max [2]float64, data interface{}, item bool) (dist float64),
) {
	return func(min, max [2]float64, data interface{}, item bool) (dist float64) {
		if item && itemDist != nil {
			return itemDist(min, max, data)
		}
		return BoxDistCalc(targetMin, targetMax, min, max, wrapX)
	}
}

func mmin(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func mmax(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

// BoxDistCalc returns the distance from rectangle A to rectangle B. When wrapX
// is provided, the operation does a cylinder wrapping of the X value to allow
// for antimeridian calculations.
func BoxDistCalc(aMin, aMax, bMin, bMax [2]float64, wrapX bool) float64 {
	var dist float64
	var squared float64

	// X
	squared = mmax(aMin[0], bMin[0]) - mmin(aMax[0], bMax[0])
	if wrapX {
		squaredLeft := mmax(aMin[0]-360, bMin[0]) - mmin(aMax[0]-360, bMax[0])
		squaredRight := mmax(aMin[0]+360, bMin[0]) - mmin(aMax[0]+360, bMax[0])
		squared = mmin(squared, mmin(squaredLeft, squaredRight))
	}
	if squared > 0 {
		dist += squared * squared
	}

	// Y
	squared = mmax(aMin[1], bMin[1]) - mmin(aMax[1], bMax[1])
	if squared > 0 {
		dist += squared * squared
	}

	return dist
}

// contains return struct when b is fully contained inside of n
func intersects(aMin, aMax, bMin, bMax [2]float64) bool {
	if bMin[0] > aMax[0] || bMax[0] < aMin[0] {
		return false
	}
	if bMin[1] > aMax[1] || bMax[1] < aMin[1] {
		return false
	}
	return true
}
