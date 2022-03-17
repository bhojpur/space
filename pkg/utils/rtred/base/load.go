package base

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

// Load bulk load items into the R-tree.
func (tr *RTree) Load(mins, maxs [][]float64, items []interface{}) {
	if len(items) < tr.minEntries {
		for i := 0; i < len(items); i++ {
			tr.Insert(mins[i], maxs[i], items[i])
		}
		return
	}

	// prefill the items
	fitems := make([]*treeNode, len(items))
	for i := 0; i < len(items); i++ {
		item := &treeItem{min: mins[i], max: maxs[i], item: items[i]}
		fitems[i] = item.unsafeNode()
	}

	// following equations are defined in the paper describing OMT
	N := len(fitems)
	M := tr.maxEntries
	h := int(math.Ceil(math.Log(float64(N)) / math.Log(float64(M))))
	Nsubtree := int(math.Pow(float64(M), float64(h-1)))
	S := int(math.Ceil(math.Sqrt(float64(N) / float64(Nsubtree))))

	// sort by the initial axis
	axis := 0
	sortByAxis(fitems, axis)

	// build the root node. it's split differently from the subtrees.
	children := make([]*treeNode, 0, S)
	for i := 0; i < S; i++ {
		var part []*treeNode
		if i == S-1 {
			// last split
			part = fitems[len(fitems)/S*i:]
		} else {
			part = fitems[len(fitems)/S*i : len(fitems)/S*(i+1)]
		}
		children = append(children, tr.omt(part, h-1, axis+1))
	}

	node := tr.createNode(children)
	node.leaf = false
	node.height = h
	tr.calcBBox(node)

	if tr.data.count == 0 {
		// save as is if tree is empty
		tr.data = node
	} else if tr.data.height == node.height {
		// split root if trees have the same height
		tr.splitRoot(tr.data, node)
	} else {
		if tr.data.height < node.height {
			// swap trees if inserted one is bigger
			tr.data, node = node, tr.data
		}

		// insert the small tree into the large tree at appropriate level
		tr.insert(node, nil, tr.data.height-node.height-1, true)
	}
}

func (tr *RTree) omt(fitems []*treeNode, h, axis int) *treeNode {
	if len(fitems) <= tr.maxEntries {
		// reached leaf level; return leaf
		children := make([]*treeNode, len(fitems))
		copy(children, fitems)
		node := tr.createNode(children)
		node.height = h
		tr.calcBBox(node)
		return node
	}

	// sort the items on a different axis than the previous level.
	sortByAxis(fitems, axis%tr.dims)
	children := make([]*treeNode, 0, tr.maxEntries)
	partsz := len(fitems) / tr.maxEntries
	for i := 0; i < tr.maxEntries; i++ {
		var part []*treeNode
		if i == tr.maxEntries-1 {
			// last part
			part = fitems[partsz*i:]
		} else {
			part = fitems[partsz*i : partsz*(i+1)]
		}
		children = append(children, tr.omt(part, h-1, axis+1))
	}
	node := tr.createNode(children)
	node.height = h
	node.leaf = false
	tr.calcBBox(node)
	return node
}
