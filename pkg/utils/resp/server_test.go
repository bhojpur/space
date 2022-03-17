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
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	// Use the example server in example_test
	go func() {
		ExampleServer()
	}()
	if os.Getenv("WAIT_ON_TEST_SERVER") == "1" {
		select {}
	}
	time.Sleep(time.Millisecond * 50)

	n := 75

	// Open N connections and do a bunch of stuff.
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			nconn, err := net.Dial("tcp", ":6380")
			if err != nil {
				t.Fatal(err)
			}
			defer nconn.Close()
			conn := NewConn(nconn)

			// PING
			if err := conn.WriteMultiBulk("PING"); err != nil {
				t.Fatal(err)
			}
			val, _, err := conn.ReadValue()
			if err != nil {
				t.Fatal(err)
			}
			if val.String() != "PONG" {
				t.Fatalf("expecting 'PONG', got '%s'", val)
			}

			key := fmt.Sprintf("key:%d", i)

			// SET
			if err := conn.WriteMultiBulk("SET", key, 123.4); err != nil {
				t.Fatal(err)
			}
			val, _, err = conn.ReadValue()
			if err != nil {
				t.Fatal(err)
			}
			if val.String() != "OK" {
				t.Fatalf("expecting 'OK', got '%s'", val)
			}

			// GET
			if err := conn.WriteMultiBulk("GET", key); err != nil {
				t.Fatal(err)
			}
			val, _, err = conn.ReadValue()
			if err != nil {
				t.Fatal(err)
			}
			if val.Float() != 123.4 {
				t.Fatalf("expecting '123.4', got '%s'", val)
			}

			// QUIT
			if err := conn.WriteMultiBulk("QUIT"); err != nil {
				t.Fatal(err)
			}
			val, _, err = conn.ReadValue()
			if err != nil {
				t.Fatal(err)
			}
			if val.String() != "OK" {
				t.Fatalf("expecting 'OK', got '%s'", val)
			}

		}(i)
	}
	wg.Wait()
}
