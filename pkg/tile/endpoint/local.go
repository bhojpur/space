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

// LocalPublisher is used to publish local notifcations
type LocalPublisher interface {
	Publish(channel string, message ...string) int
}

// LocalConn is an endpoint connection
type LocalConn struct {
	ep        Endpoint
	publisher LocalPublisher
}

func newLocalConn(ep Endpoint, publisher LocalPublisher) *LocalConn {
	return &LocalConn{
		ep:        ep,
		publisher: publisher,
	}
}

// Expired returns true if the connection has expired
func (conn *LocalConn) Expired() bool {
	return false
}

// Send sends a message
func (conn *LocalConn) Send(msg string) error {
	conn.publisher.Publish(conn.ep.Local.Channel, msg)
	return nil
}
