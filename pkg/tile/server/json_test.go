package server

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
	"encoding/json"
	"testing"
)

func BenchmarkJSONString(t *testing.B) {
	var s = "the need for mead"
	for i := 0; i < t.N; i++ {
		jsonString(s)
	}
}

func BenchmarkJSONMarshal(t *testing.B) {
	var s = "the need for mead"
	for i := 0; i < t.N; i++ {
		json.Marshal(s)
	}
}

func TestIsJsonNumber(t *testing.T) {
	test := func(expected bool, val string) {
		actual := isJSONNumber(val)
		if expected != actual {
			t.Fatalf("Expected %t == isJsonNumber(\"%s\") but was %t", expected, val, actual)
		}
	}
	test(false, "")
	test(false, "-")
	test(false, "foo")
	test(false, "0123")
	test(false, "1.")
	test(false, "1.0e")
	test(false, "1.0e-")
	test(false, "1.0E10NaN")
	test(false, "1.0ENaN")
	test(true, "-1")
	test(true, "0")
	test(true, "0.0")
	test(true, "42")
	test(true, "1.0E10")
	test(true, "1.0e10")
	test(true, "1E+5")
	test(true, "1E-10")
}
