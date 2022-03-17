package log

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
	"io/ioutil"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLog(t *testing.T) {
	f := &bytes.Buffer{}
	LogJSON = false
	SetOutput(f)
	Printf("hello %v", "everyone")
	if !strings.HasSuffix(f.String(), "hello everyone\n") {
		t.Fatal("fail")
	}
}

func TestLogJSON(t *testing.T) {

	LogJSON = true
	Build("")

	type tcase struct {
		level  int
		format string
		args   string
		ops    func(...interface{})
		fops   func(string, ...interface{})
		expMsg string
		expLvl zapcore.Level
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(zap.DebugLevel)
			Set(zap.New(observedZapCore).Sugar())
			Level = tc.level

			if tc.format != "" {
				tc.fops(tc.format, tc.args)
			} else {
				tc.ops(tc.args)
			}

			if observedLogs.Len() < 1 {
				t.Fatal("fail")
			}

			allLogs := observedLogs.All()

			if allLogs[0].Message != tc.expMsg {
				t.Fatal("fail")
			}

			if allLogs[0].Level != tc.expLvl {
				t.Fatal("fail")
			}
		}
	}

	tests := map[string]tcase{
		"Print": {
			level: 1,
			args:  "Print json logger",
			ops: func(args ...interface{}) {
				Print(args...)
			},
			expMsg: "Print json logger",
			expLvl: zapcore.InfoLevel,
		},
		"Printf": {
			level:  1,
			format: "Printf json %v",
			args:   "logger",
			fops: func(format string, args ...interface{}) {
				Printf(format, args...)
			},
			expMsg: "Printf json logger",
			expLvl: zapcore.InfoLevel,
		},
		"Info": {
			level: 1,
			args:  "Info json logger",
			ops: func(args ...interface{}) {
				Info(args...)
			},
			expMsg: "Info json logger",
			expLvl: zapcore.InfoLevel,
		},
		"Infof": {
			level:  1,
			format: "Infof json %v",
			args:   "logger",
			fops: func(format string, args ...interface{}) {
				Infof(format, args...)
			},
			expMsg: "Infof json logger",
			expLvl: zapcore.InfoLevel,
		},
		"Debug": {
			level: 3,
			args:  "Debug json logger",
			ops: func(args ...interface{}) {
				Debug(args...)
			},
			expMsg: "Debug json logger",
			expLvl: zapcore.DebugLevel,
		},
		"Debugf": {
			level:  3,
			format: "Debugf json %v",
			args:   "logger",
			fops: func(format string, args ...interface{}) {
				Debugf(format, args...)
			},
			expMsg: "Debugf json logger",
			expLvl: zapcore.DebugLevel,
		},
		"Warn": {
			level: 2,
			args:  "Warn json logger",
			ops: func(args ...interface{}) {
				Warn(args...)
			},
			expMsg: "Warn json logger",
			expLvl: zapcore.WarnLevel,
		},
		"Warnf": {
			level:  2,
			format: "Warnf json %v",
			args:   "logger",
			fops: func(format string, args ...interface{}) {
				Warnf(format, args...)
			},
			expMsg: "Warnf json logger",
			expLvl: zapcore.WarnLevel,
		},
		"Error": {
			level: 1,
			args:  "Error json logger",
			ops: func(args ...interface{}) {
				Error(args...)
			},
			expMsg: "Error json logger",
			expLvl: zapcore.ErrorLevel,
		},
		"Errorf": {
			level:  1,
			format: "Errorf json %v",
			args:   "logger",
			fops: func(format string, args ...interface{}) {
				Errorf(format, args...)
			},
			expMsg: "Errorf json logger",
			expLvl: zapcore.ErrorLevel,
		},
		"Http": {
			level: 1,
			args:  "Http json logger",
			ops: func(args ...interface{}) {
				HTTP(args...)
			},
			expMsg: "Http json logger",
			expLvl: zapcore.InfoLevel,
		},
		"Httpf": {
			level:  1,
			format: "Httpf json %v",
			args:   "logger",
			fops: func(format string, args ...interface{}) {
				HTTPf(format, args...)
			},
			expMsg: "Httpf json logger",
			expLvl: zapcore.InfoLevel,
		},
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}

func BenchmarkLogPrintf(t *testing.B) {
	LogJSON = false
	Level = 1
	SetOutput(ioutil.Discard)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Printf("X %s", "Y")
	}
}

func BenchmarkLogJSONPrintf(t *testing.B) {
	LogJSON = true
	Level = 1

	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)

	logger := zap.New(
		zapcore.NewCore(
			enc,
			zapcore.AddSync(ioutil.Discard),
			zap.DebugLevel,
		)).Sugar()

	Set(logger)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Printf("X %s", "Y")
	}
}
