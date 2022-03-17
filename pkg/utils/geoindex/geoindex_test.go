package geoindex

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
	"math/rand"
	"testing"
	"time"

	"github.com/bhojpur/space/pkg/utils/geoindex/internal"
)

func init() {
	seed := time.Now().UnixNano()
	println("seed:", seed)
	rand.Seed(seed)
}

func TestGeoIndex(t *testing.T) {
	t.Run("BenchVarious", func(t *testing.T) {
		Tests.TestBenchVarious(t, &internal.RTree{}, 100000)
	})
	t.Run("RandomRects", func(t *testing.T) {
		Tests.TestRandomRects(t, &internal.RTree{}, 10000)
	})
	t.Run("RandomPoints", func(t *testing.T) {
		Tests.TestRandomPoints(t, &internal.RTree{}, 10000)
	})
	t.Run("ZeroPoints", func(t *testing.T) {
		Tests.TestZeroPoints(t, &internal.RTree{})
	})
	t.Run("CitiesSVG", func(t *testing.T) {
		Tests.TestCitiesSVG(t, &internal.RTree{})
	})
}

func BenchmarkRandomInsert(b *testing.B) {
	Tests.BenchmarkRandomInsert(b, &internal.RTree{})
}
