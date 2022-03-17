package endpoint

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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	httpExpiresAfter       = time.Second * 30
	httpRequestTimeout     = time.Second * 5
	httpMaxIdleConnections = 20
)

// HTTPConn is an endpoint connection
type HTTPConn struct {
	ep     Endpoint
	client *http.Client
}

func newHTTPConn(ep Endpoint) *HTTPConn {
	return &HTTPConn{
		ep: ep,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: httpMaxIdleConnections,
				IdleConnTimeout:     httpExpiresAfter,
			},
			Timeout: httpRequestTimeout,
		},
	}
}

// Expired returns true if the connection has expired
func (conn *HTTPConn) Expired() bool {
	return false
}

// Send sends a message
func (conn *HTTPConn) Send(msg string) error {
	req, err := http.NewRequest("POST", conn.ep.Original, bytes.NewBufferString(msg))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := conn.client.Do(req)
	if err != nil {
		return err
	}
	// close the connection to reuse it
	defer resp.Body.Close()
	// discard response
	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		return err
	}
	// Only allow responses with status code 200, 201, and 202
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("invalid status: %s", resp.Status)
	}
	return nil
}
