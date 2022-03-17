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
)

func TestQueue(t *testing.T) {
	var q queue

	q.push(qnode{
		dist: 2,
	})
	q.push(qnode{
		dist: 1,
	})
	q.push(qnode{
		dist: 5,
	})
	q.push(qnode{
		dist: 3,
	})
	q.push(qnode{
		dist: 4,
	})

	lastDist := float64(-1)
	for i := 0; i < 3; i++ {
		node, ok := q.pop()
		if !ok {
			t.Fatal("queue was empty")
		}
		if node.dist < lastDist {
			t.Fatal("queue was out of order")
		}
	}

	if len(q) != 2 {
		t.Fatal("queue was wrong size")
	}

	capBeforeInserts := cap(q)
	q.push(qnode{
		dist: 1,
	})
	q.push(qnode{
		dist: 10,
	})
	q.push(qnode{
		dist: 11,
	})

	if cap(q) != capBeforeInserts {
		t.Fatal("queue did not reuse space")
	}

	lastDist = -1
	for i := 0; i < 5; i++ {
		node, ok := q.pop()
		if !ok {
			t.Fatal("queue was empty")
		}
		if node.dist < lastDist {
			t.Fatal("queue was out of order")
		}
	}

	_, ok := q.pop()
	if ok {
		t.Fatal("queue was not empty")
	}
}

func BenchmarkQueue(b *testing.B) {
	var q queue

	for i := 0; i < b.N; i++ {
		r := rand.Float64()
		if r < 0.5 {
			q.push(qnode{dist: r})
		} else {
			q.pop()
		}
	}
}
