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
	"sync/atomic"
)

type aint struct{ v uintptr }

func (a *aint) add(d int) int {
	if d < 0 {
		return int(atomic.AddUintptr(&a.v, ^uintptr((d*-1)-1)))
	}
	return int(atomic.AddUintptr(&a.v, uintptr(d)))
}
func (a *aint) get() int {
	return int(atomic.LoadUintptr(&a.v))
}
func (a *aint) set(i int) int {
	return int(atomic.SwapUintptr(&a.v, uintptr(i)))
}

type abool struct{ v uint32 }

func (a *abool) on() bool {
	return atomic.LoadUint32(&a.v) != 0
}
func (a *abool) set(t bool) bool {
	if t {
		return atomic.SwapUint32(&a.v, 1) != 0
	}
	return atomic.SwapUint32(&a.v, 0) != 0
}
