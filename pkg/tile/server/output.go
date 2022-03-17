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
	"strings"
	"time"

	"github.com/bhojpur/space/pkg/utils/resp"
)

func (s *Server) cmdOutput(msg *Message) (res resp.Value, err error) {
	start := time.Now()
	vs := msg.Args[1:]
	var arg string
	var ok bool

	if len(vs) != 0 {
		if _, arg, ok = tokenval(vs); !ok || arg == "" {
			return NOMessage, errInvalidNumberOfArguments
		}
		// Setting the original message output type will be picked up by the
		// server prior to the next command being executed.
		switch strings.ToLower(arg) {
		default:
			return NOMessage, errInvalidArgument(arg)
		case "json":
			msg.OutputType = JSON
		case "resp":
			msg.OutputType = RESP
		}
		return OKMessage(msg, start), nil
	}
	// return the output
	switch msg.OutputType {
	default:
		return NOMessage, nil
	case JSON:
		return resp.StringValue(`{"ok":true,"output":"json","elapsed":` + time.Since(start).String() + `}`), nil
	case RESP:
		return resp.StringValue("resp"), nil
	}
}
