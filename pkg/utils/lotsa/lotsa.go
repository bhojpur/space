package lotsa

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
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

// Output is used to print elased time and ops/sec
var Output io.Writer

// MemUsage is used to output the memory usage
var MemUsage bool

// Ops executed a number of operations over a multiple goroutines.
// count is the number of operations.
// threads is the number goroutines.
// op is the operation function
func Ops(count, threads int, op func(i, thread int)) {
	var start time.Time
	var wg sync.WaitGroup
	wg.Add(threads)
	var ms1 runtime.MemStats
	output := Output
	if output != nil {
		if MemUsage {
			runtime.GC()
			runtime.ReadMemStats(&ms1)
		}
		start = time.Now()
	}
	for i := 0; i < threads; i++ {
		s, e := count/threads*i, count/threads*(i+1)
		if i == threads-1 {
			e = count
		}
		go func(i, s, e int) {
			defer wg.Done()
			for j := s; j < e; j++ {
				op(j, i)
			}
		}(i, s, e)
	}
	wg.Wait()

	if output != nil {
		dur := time.Since(start)
		var alloc uint64
		if MemUsage {
			runtime.GC()
			var ms2 runtime.MemStats
			runtime.ReadMemStats(&ms2)
			if ms1.HeapAlloc > ms2.HeapAlloc {
				alloc = 0
			} else {
				alloc = ms2.HeapAlloc - ms1.HeapAlloc
			}
		}
		WriteOutput(output, count, threads, dur, alloc)
	}
}

func commaize(n int) string {
	s1, s2 := fmt.Sprintf("%d", n), ""
	for i, j := len(s1)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j%3 == 0 && j != 0 {
			s2 = "," + s2
		}
		s2 = string(s1[i]) + s2
	}
	return s2
}

func memstr(alloc uint64) string {
	switch {
	case alloc <= 1024:
		return fmt.Sprintf("%d bytes", alloc)
	case alloc <= 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(alloc)/1024)
	case alloc <= 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(alloc)/1024/1024)
	default:
		return fmt.Sprintf("%.1f GB", float64(alloc)/1024/1024/1024)
	}
}

// WriteOutput writes an output line to the specified writer
func WriteOutput(w io.Writer, count, threads int, elapsed time.Duration, alloc uint64) {
	var ss string
	if threads != 1 {
		ss = fmt.Sprintf("over %d threads ", threads)
	}
	var nsop int
	if count > 0 {
		nsop = int(elapsed / time.Duration(count))
	}
	var allocstr string
	if alloc > 0 {
		var bops uint64
		if count > 0 {
			bops = alloc / uint64(count)
		}
		allocstr = fmt.Sprintf(", %s, %d bytes/op", memstr(alloc), bops)
	}
	fmt.Fprintf(w, "%s ops %sin %.0fms, %s/sec, %d ns/op%s\n",
		commaize(count), ss, elapsed.Seconds()*1000,
		commaize(int(float64(count)/elapsed.Seconds())),
		nsop, allocstr,
	)
}
