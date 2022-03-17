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
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const amqpExpiresAfter = time.Second * 30

// AMQPConn is an endpoint connection
type AMQPConn struct {
	mu      sync.Mutex
	ep      Endpoint
	conn    *amqp.Connection
	channel *amqp.Channel
	ex      bool
	t       time.Time
}

// Expired returns true if the connection has expired
func (conn *AMQPConn) Expired() bool {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.ex {
		if time.Since(conn.t) > amqpExpiresAfter {
			conn.ex = true
			conn.close()
		}
	}
	return conn.ex
}

func (conn *AMQPConn) close() {
	if conn.conn != nil {
		conn.conn.Close()
		conn.conn = nil
		conn.channel = nil
	}
}

// Send sends a message
func (conn *AMQPConn) Send(msg string) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.ex {
		return errExpired
	}
	conn.t = time.Now()

	if conn.conn == nil {
		prefix := "amqp://"
		if conn.ep.AMQP.SSL {
			prefix = "amqps://"
		}

		var cfg amqp.Config
		cfg.Dial = func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Second)
		}
		c, err := amqp.DialConfig(fmt.Sprintf("%s%s", prefix, conn.ep.AMQP.URI), cfg)

		if err != nil {
			return err
		}

		channel, err := c.Channel()
		if err != nil {
			return err
		}

		// Declare new exchange
		if err := channel.ExchangeDeclare(
			conn.ep.AMQP.QueueName,
			conn.ep.AMQP.Type,
			conn.ep.AMQP.Durable,
			conn.ep.AMQP.AutoDelete,
			conn.ep.AMQP.Internal,
			conn.ep.AMQP.NoWait,
			nil,
		); err != nil {
			return err
		}

		// Create queue if queue don't exists
		if _, err := channel.QueueDeclare(
			conn.ep.AMQP.QueueName,
			conn.ep.AMQP.Durable,
			conn.ep.AMQP.AutoDelete,
			false,
			conn.ep.AMQP.NoWait,
			nil,
		); err != nil {
			return err
		}

		// Binding exchange to queue
		if err := channel.QueueBind(
			conn.ep.AMQP.QueueName,
			conn.ep.AMQP.RouteKey,
			conn.ep.AMQP.QueueName,
			conn.ep.AMQP.NoWait,
			nil,
		); err != nil {
			return err
		}

		conn.conn = c
		conn.channel = channel
	}

	return conn.channel.Publish(
		conn.ep.AMQP.QueueName,
		conn.ep.AMQP.RouteKey,
		conn.ep.AMQP.Mandatory,
		conn.ep.AMQP.Immediate,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(msg),
			DeliveryMode:    conn.ep.AMQP.DeliveryMode,
			Priority:        conn.ep.AMQP.Priority,
		},
	)
}

func newAMQPConn(ep Endpoint) *AMQPConn {
	return &AMQPConn{
		ep: ep,
		t:  time.Now(),
	}
}
