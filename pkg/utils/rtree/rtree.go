package rtree

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

type RTree struct {
	base Generic[any]
}

// Insert an item into the structure
func (tr *RTree) Insert(min, max [2]float64, data interface{}) {
	tr.base.Insert(min, max, data)
}

// Delete an item from the structure
func (tr *RTree) Delete(min, max [2]float64, data interface{}) {
	tr.base.Delete(min, max, data)
}

// Replace an item in the structure. This is effectively just a Delete
// followed by an Insert. But for some structures it may be possible to
// optimize the operation to avoid multiple passes
func (tr *RTree) Replace(
	oldMin, oldMax [2]float64, oldData interface{},
	newMin, newMax [2]float64, newData interface{},
) {
	tr.base.Replace(
		oldMin, oldMax, oldData,
		newMin, newMax, newData,
	)
}

// Search the structure for items that intersects the rect param
func (tr *RTree) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, data interface{}) bool,
) {
	tr.base.Search(min, max, iter)
}

// Scan iterates through all data in tree in no specified order.
func (tr *RTree) Scan(iter func(min, max [2]float64, data interface{}) bool) {
	tr.base.Scan(iter)
}

// Len returns the number of items in tree
func (tr *RTree) Len() int {
	return tr.base.Len()
}

// Bounds returns the minimum bounding box
func (tr *RTree) Bounds() (min, max [2]float64) {
	return tr.base.Bounds()
}

// Children returns all children for parent node. If parent node is nil
// then the root nodes should be returned.
// The reuse buffer is an empty length slice that can optionally be used
// to avoid extra allocations.
func (tr *RTree) Children(parent interface{}, reuse []child.Child) (children []child.Child) {
	return tr.base.Children(parent, reuse)
}
