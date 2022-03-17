package resp

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
	"os"
	"testing"
	"time"
)

func TestAOF(t *testing.T) {
	os.RemoveAll("aof.tmp")
	if err := os.MkdirAll("aof.tmp", 0700); err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("aof.tmp")
	}()

	if _, err := OpenAOF("aof.tmp"); err == nil {
		t.Fatal("expecting error, got nil")
	}

	f, err := OpenAOF("aof.tmp/aof")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		if err := f.Close(); err == nil || err.Error() != "closed" {
			t.Fatalf("expected 'closed', got '%v'", err)
		}
	}()
	// Test Setting Sync Policies
	f.SetSyncPolicy(Never)
	sps := fmt.Sprintf("%s %s %s %s", SyncPolicy(99), Never, Always, EverySecond)
	if sps != "unknown never always every second" {
		t.Fatalf("expected '%s', got '%s'", "unknown never always every second", sps)
	}
	f.SetSyncPolicy(99)
	f.SetSyncPolicy(Never)
	f.SetSyncPolicy(Always)
	f.SetSyncPolicy(EverySecond)
	f.SetSyncPolicy(EverySecond)
	for i := 0; i < 12345; i++ {
		if err := f.Append(StringValue(fmt.Sprintf("hello world #%d\n", i))); err != nil {
			t.Fatal(err)
		}
	}
	i := 0
	if err := f.Scan(func(v Value) {
		s := v.String()
		e := fmt.Sprintf("hello world #%d\n", i)
		if s != e {
			t.Fatalf("#%d is '%s', expect '%s'", i, s, e)
		}
		i++
	}); err != nil {
		t.Fatal(err)
	}
	f.Close()
	f.Close() // Test closing twice
	f, err = OpenAOF("aof.tmp/aof")
	if err != nil {
		t.Fatal(err)
	}
	c := i
	for i := c; i < c+12345; i++ {
		if err := f.Append(StringValue(fmt.Sprintf("hello world #%d\n", i))); err != nil {
			t.Fatal(err)
		}
	}
	i = 0
	if err := f.Scan(func(v Value) {
		s := v.String()
		e := fmt.Sprintf("hello world #%d\n", i)
		if s != e {
			t.Fatalf("#%d is '%s', expect '%s'", i, s, e)
		}
		i++
	}); err != nil {
		t.Fatal(err)
	}

	var multi []Value
	for i := 0; i < 50; i++ {
		multi = append(multi, StringValue(fmt.Sprintf("hello multi world #%d\n", i)))
	}
	if err := f.AppendMulti(multi); err != nil {
		t.Fatal(err)
	}

	skip := i
	i = 0
	j := 0
	if err := f.Scan(func(v Value) {
		if i >= skip {
			s := v.String()
			e := fmt.Sprintf("hello multi world #%d\n", j)
			if s != e {
				t.Fatalf("#%d is '%s', expect '%s'", j, s, e)
			}
			j++
		}
		i++
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 10)
}
