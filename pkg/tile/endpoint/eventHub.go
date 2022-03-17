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
	"time"

	"github.com/bhojpur/space/pkg/utils/gjson"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

const ()

// HTTPConn is an endpoint connection
type EvenHubConn struct {
	ep Endpoint
}

func newEventHubConn(ep Endpoint) *EvenHubConn {
	return &EvenHubConn{
		ep: ep,
	}
}

// Expired returns true if the connection has expired
func (conn *EvenHubConn) Expired() bool {
	return false
}

// Send sends a message
func (conn *EvenHubConn) Send(msg string) error {
	hub, err := eventhub.NewHubFromConnectionString(conn.ep.EventHub.ConnectionString)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// parse json again to get out info for our kafka key
	key := gjson.Get(msg, "key")
	id := gjson.Get(msg, "id")
	keyValue := fmt.Sprintf("%s-%s", key.String(), id.String())

	evtHubMsg := eventhub.NewEventFromString(msg)
	evtHubMsg.PartitionKey = &keyValue
	err = hub.Send(ctx, evtHubMsg)
	if err != nil {
		return err
	}

	return nil
}
