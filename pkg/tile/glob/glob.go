package glob

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

import "strings"

// Glob structure for simple string matching
type Glob struct {
	Pattern string
	Desc    bool
	Limits  []string
	IsGlob  bool
}

// Match returns true when string matches pattern. Returns an error when the
// pattern is invalid.
func Match(pattern, str string) (matched bool, err error) {
	return wildcardMatch(pattern, str)
}

// IsGlob returns true when the pattern is a valid glob
func IsGlob(pattern string) bool {
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '[', '*', '?':
			_, err := Match(pattern, "whatever")
			return err == nil
		}
	}
	return false
}

// Parse returns a glob structure from the pattern.
func Parse(pattern string, desc bool) *Glob {
	g := &Glob{Pattern: pattern, Desc: desc, Limits: []string{"", ""}}
	if strings.HasPrefix(pattern, "*") {
		g.IsGlob = true
		return g
	}
	if pattern == "" {
		g.IsGlob = false
		return g
	}
	n := 0
	isGlob := false
outer:
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '[', '*', '?':
			_, err := Match(pattern, "whatever")
			if err == nil {
				isGlob = true
			}
			break outer
		}
		n++
	}
	if n == 0 {
		g.Limits = []string{pattern, pattern}
		g.IsGlob = false
		return g
	}
	var a, b string
	if desc {
		a = pattern[:n]
		b = a
		if b[n-1] == 0x00 {
			for len(b) > 0 && b[len(b)-1] == 0x00 {
				if len(b) > 1 {
					if b[len(b)-2] == 0x00 {
						b = b[:len(b)-1]
					} else {
						b = string(append([]byte(b[:len(b)-2]), b[len(b)-2]-1, 0xFF))
					}
				} else {
					b = ""
				}
			}
		} else {
			b = string(append([]byte(b[:n-1]), b[n-1]-1))
		}
		if a[n-1] == 0xFF {
			a = string(append([]byte(a), 0x00))
		} else {
			a = string(append([]byte(a[:n-1]), a[n-1]+1))
		}
	} else {
		a = pattern[:n]
		if a[n-1] == 0xFF {
			b = string(append([]byte(a), 0x00))
		} else {
			b = string(append([]byte(a[:n-1]), a[n-1]+1))
		}
	}
	g.Limits = []string{a, b}
	g.IsGlob = isGlob
	return g
}
