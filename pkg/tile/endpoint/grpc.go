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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bhojpur/space/pkg/tile/hservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const grpcExpiresAfter = time.Second * 30

// GRPCConn is an endpoint connection
type GRPCConn struct {
	mu    sync.Mutex
	ep    Endpoint
	ex    bool
	t     time.Time
	conn  *grpc.ClientConn
	sconn hservice.HookServiceClient
}

func newGRPCConn(ep Endpoint) *GRPCConn {
	return &GRPCConn{
		ep: ep,
		t:  time.Now(),
	}
}

// Expired returns true if the connection has expired
func (conn *GRPCConn) Expired() bool {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.ex {
		if time.Since(conn.t) > grpcExpiresAfter {
			if conn.conn != nil {
				conn.close()
			}
			conn.ex = true
		}
	}
	return conn.ex
}
func (conn *GRPCConn) close() {
	if conn.conn != nil {
		conn.conn.Close()
		conn.conn = nil
	}
}

// Send sends a message
func (conn *GRPCConn) Send(msg string) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.ex {
		return errExpired
	}
	conn.t = time.Now()
	if conn.conn == nil {
		addr := fmt.Sprintf("%s:%d", conn.ep.GRPC.Host, conn.ep.GRPC.Port)
		var err error
		conn.conn, err = grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			conn.close()
			return err
		}
		conn.sconn = hservice.NewHookServiceClient(conn.conn)
	}
	r, err := conn.sconn.Send(context.Background(), &hservice.MessageRequest{Value: msg})
	if err != nil {
		conn.close()
		return err
	}
	if !r.Ok {
		conn.close()
		return errors.New("invalid grpc reply")
	}
	return nil
}
