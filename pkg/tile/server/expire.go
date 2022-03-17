package server

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
	"time"

	"github.com/bhojpur/space/pkg/tile/log"
)

const bgExpireDelay = time.Second / 10

// backgroundExpiring deletes expired items from the database.
// It's executes every 1/10 of a second.
func (s *Server) backgroundExpiring() {
	for {
		if s.stopServer.on() {
			return
		}
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			now := time.Now()
			s.backgroundExpireObjects(now)
			s.backgroundExpireHooks(now)
		}()
		time.Sleep(bgExpireDelay)
	}
}

func (s *Server) backgroundExpireObjects(now time.Time) {
	nano := now.UnixNano()
	var ids []string
	var msgs []*Message
	s.cols.Ascend(nil, func(v interface{}) bool {
		col := v.(*collectionKeyContainer)
		ids = col.col.Expired(nano, ids[:0])
		for _, id := range ids {
			msgs = append(msgs, &Message{
				Args: []string{"del", col.key, id},
			})
		}
		return true
	})
	for _, msg := range msgs {
		_, d, err := s.cmdDel(msg)
		if err != nil {
			log.Fatal(err)
		}
		if err := s.writeAOF(msg.Args, &d); err != nil {
			log.Fatal(err)
		}
	}
	if len(msgs) > 0 {
		log.Debugf("Expired %d objects\n", len(msgs))
	}

}

func (s *Server) backgroundExpireHooks(now time.Time) {
	var msgs []*Message
	s.hookExpires.Ascend(nil, func(v interface{}) bool {
		h := v.(*Hook)
		if h.expires.After(now) {
			return false
		}
		msg := &Message{}
		if h.channel {
			msg.Args = []string{"delchan", h.Name}
		} else {
			msg.Args = []string{"delhook", h.Name}
		}
		msgs = append(msgs, msg)
		return true
	})

	for _, msg := range msgs {
		_, d, err := s.cmdDelHook(msg)
		if err != nil {
			log.Fatal(err)
		}
		if err := s.writeAOF(msg.Args, &d); err != nil {
			log.Fatal(err)
		}
	}
	if len(msgs) > 0 {
		log.Debugf("Expired %d hooks\n", len(msgs))
	}
}
