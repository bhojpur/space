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

import "github.com/bhojpur/space/pkg/utils/geoindex/child"

// Priority Queue ordered by dist (smallest to largest)

type qnode struct {
	dist  float64
	child child.Child
}

type queue []qnode

func (q *queue) push(node qnode) {
	*q = append(*q, node)
	nodes := *q
	i := len(nodes) - 1
	parent := (i - 1) / 2
	for ; i != 0 && nodes[parent].dist > nodes[i].dist; parent = (i - 1) / 2 {
		nodes[parent], nodes[i] = nodes[i], nodes[parent]
		i = parent
	}
}

func (q *queue) pop() (qnode, bool) {
	nodes := *q
	if len(nodes) == 0 {
		return qnode{}, false
	}
	var n qnode
	n, nodes[0] = nodes[0], nodes[len(*q)-1]
	nodes = nodes[:len(nodes)-1]
	*q = nodes

	i := 0
	for {
		smallest := i
		left := i*2 + 1
		right := i*2 + 2
		if left < len(nodes) && nodes[left].dist <= nodes[smallest].dist {
			smallest = left
		}
		if right < len(nodes) && nodes[right].dist <= nodes[smallest].dist {
			smallest = right
		}
		if smallest == i {
			break
		}
		nodes[smallest], nodes[i] = nodes[i], nodes[smallest]
		i = smallest
	}
	return n, true
}
