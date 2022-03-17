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
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const pubsubExpiresAfter = time.Second * 30

// SQSConn is an endpoint connection
type PubSubConn struct {
	mu    sync.Mutex
	ep    Endpoint
	svc   *pubsub.Client
	topic *pubsub.Topic
	ex    bool
	t     time.Time
}

func (conn *PubSubConn) close() {
	if conn.svc != nil {
		conn.svc.Close()
		conn.svc = nil
	}
}

// Send sends a message
func (conn *PubSubConn) Send(msg string) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.ex {
		return errExpired
	}

	ctx := context.Background()

	conn.t = time.Now()

	if conn.svc == nil {
		var creds option.ClientOption
		var svc *pubsub.Client
		var err error
		credPath := conn.ep.PubSub.CredPath

		if credPath != "" {
			creds = option.WithCredentialsFile(credPath)
			svc, err = pubsub.NewClient(ctx, conn.ep.PubSub.Project, creds)
		} else {
			svc, err = pubsub.NewClient(ctx, conn.ep.PubSub.Project)
		}

		if err != nil {
			fmt.Println(err)
			return err
		}

		topic := svc.Topic(conn.ep.PubSub.Topic)

		conn.svc = svc
		conn.topic = topic
	}

	// Send message
	res := conn.topic.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	_, err := res.Get(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (conn *PubSubConn) Expired() bool {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.ex {
		if time.Since(conn.t) > pubsubExpiresAfter {
			conn.ex = true
			conn.close()
		}
	}
	return conn.ex
}

func newPubSubConn(ep Endpoint) *PubSubConn {
	return &PubSubConn{
		ep: ep,
		t:  time.Now(),
	}
}
