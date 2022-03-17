package deadline

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

import "time"

// Deadline allows for commands to expire when they run too long
type Deadline struct {
	unixNano int64
	hit      bool
}

// New returns a new deadline object
func New(dl time.Time) *Deadline {
	return &Deadline{unixNano: dl.UnixNano()}
}

// Check the deadline and panic when reached
//go:noinline
func (dl *Deadline) Check() {
	if dl == nil || dl.unixNano == 0 {
		return
	}
	if !dl.hit && time.Now().UnixNano() > dl.unixNano {
		dl.hit = true
		panic("deadline")
	}
}

// Hit returns true if the deadline has been hit
func (dl *Deadline) Hit() bool {
	return dl.hit
}

// GetDeadlineTime returns the time object for the deadline, and an
// "empty" boolean
func (dl *Deadline) GetDeadlineTime() time.Time {
	return time.Unix(0, dl.unixNano)
}
