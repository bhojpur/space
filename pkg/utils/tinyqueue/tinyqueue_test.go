package tinyqueue

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
	"reflect"
	"sort"
	"testing"
	"time"
)

type floatValue float64

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("'%v' != '%v'", a, b)
	}
}

func (a floatValue) Less(b Item) bool {
	return a < b.(floatValue)
}

var data, sorted = func() ([]Item, []Item) {
	rand.Seed(time.Now().UnixNano())
	var data []Item
	for i := 0; i < 100; i++ {
		data = append(data, floatValue(rand.Float64()*100))
	}
	sorted := make([]Item, len(data))
	copy(sorted, data)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Less(sorted[j])
	})
	return data, sorted
}()

func TestMaintainsPriorityQueue(t *testing.T) {
	q := New(nil)
	for i := 0; i < len(data); i++ {
		q.Push(data[i])
	}
	assertEqual(t, q.Peek(), sorted[0])
	var result []Item
	for q.length > 0 {
		result = append(result, q.Pop())
	}
	assertEqual(t, result, sorted)
}

func TestAcceptsDataInConstructor(t *testing.T) {
	q := New(data)
	var result []Item
	for q.length > 0 {
		result = append(result, q.Pop())
	}
	assertEqual(t, result, sorted)
}
func TestHandlesEdgeCasesWithFewElements(t *testing.T) {
	q := New(nil)
	q.Push(floatValue(2))
	q.Push(floatValue(1))
	q.Pop()
	q.Pop()
	q.Pop()
	q.Push(floatValue(2))
	q.Push(floatValue(1))
	assertEqual(t, float64(q.Pop().(floatValue)), 1.0)
	assertEqual(t, float64(q.Pop().(floatValue)), 2.0)
	assertEqual(t, q.Pop(), nil)
}
