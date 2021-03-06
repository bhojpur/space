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

import (
	"fmt"
	"strings"
	"testing"
)

func TestSegmentRaycast(t *testing.T) {
	// This is full coverage raycast test. It uses a 7x7 grid where the each
	// point is checked for a total of 49 tests per segment. There are 16
	// segments at 0º,30º,45º,90º,180º in both directions for a total of 784
	// tests.
	segs := []interface{}{
		S(1, 1, 5, 5), "A",
		S(5, 5, 1, 1), "B",
		[7]string{
			"-------",
			"-----*-",
			"++++*--",
			"+++*---",
			"++*----",
			"+*-----",
			"-------",
		},
		S(1, 5, 5, 1), "C",
		S(5, 1, 1, 5), "D",
		[7]string{
			"-------",
			"-*-----",
			"++*----",
			"+++*---",
			"++++*--",
			"+++++*-",
			"-------",
		},
		S(1, 3, 5, 3), "E",
		S(5, 3, 1, 3), "F",
		[7]string{
			"-------",
			"-------",
			"-------",
			"-*****-",
			"-------",
			"-------",
			"-------",
		},
		S(3, 5, 3, 1), "G",
		S(3, 1, 3, 5), "H",
		[7]string{
			"-------",
			"---*---",
			"+++*---",
			"+++*---",
			"+++*---",
			"+++*---",
			"-------",
		},
		S(1, 2, 5, 4), "I",
		S(5, 4, 1, 2), "J",
		[7]string{
			"-------",
			"-------",
			"-----*-",
			"+++*---",
			"+*-----",
			"-------",
			"-------",
		},
		S(1, 4, 5, 2), "K",
		S(5, 2, 1, 4), "L",
		[7]string{
			"-------",
			"-------",
			"-*-----",
			"+++*---",
			"+++++*-",
			"-------",
			"-------",
		},
		S(2, 1, 4, 5), "M",
		S(4, 5, 2, 1), "N",
		[7]string{
			"-------",
			"----*--",
			"++++---",
			"+++*---",
			"+++----",
			"++*----",
			"-------",
		},
		S(2, 5, 4, 1), "O",
		S(4, 1, 2, 5), "P",
		[7]string{
			"-------",
			"--*----",
			"+++----",
			"+++*---",
			"++++---",
			"++++*--",
			"-------",
		},
		S(3, 3, 3, 3), "Q",
		S(3, 3, 3, 3), "R",
		[7]string{
			"-------",
			"-------",
			"-------",
			"---*---",
			"-------",
			"-------",
			"-------",
		},
	}

	var ms string
	for i := 0; i < len(segs); i += 5 {
		segs := []interface{}{
			segs[i], segs[i+1], segs[i+4],
			segs[i+2], segs[i+3], segs[i+4],
		}
		for i := 0; i < len(segs); i += 3 {
			seg := segs[i].(Segment)
			label := segs[i+1].(string)
			grid := segs[i+2].([7]string)
			//
			var ngrid [7]string
			for y, sy := 0, 6; y < 7; y, sy = y+1, sy-1 {
				var nline string
				for x := 0; x < 7; x++ {
					// ch := grid[sy][x]
					pt := Point{float64(x), float64(y)}
					res := seg.Raycast(pt)
					if res.In {
						nline += "+"
					} else if res.On {
						nline += "*"
					} else {
						nline += "-"
					}
				}
				ngrid[sy] = nline
			}
			if grid != ngrid {
				ms += fmt.Sprintf("MISMATCH (%s) SEGMENT: %v\n", label, seg)
				ms += fmt.Sprintf("EXPECTED\n%s\n", strings.Join(grid[:], "\n"))
				ms += fmt.Sprintf("GOT\n%s\n", strings.Join(ngrid[:], "\n"))
			}
		}
	}
	if ms != "" {
		t.Fatalf("\n%s", ms)
	}
}
