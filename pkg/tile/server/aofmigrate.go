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
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/bhojpur/space/pkg/tile/log"
	"github.com/bhojpur/space/pkg/utils/resp"
)

var errCorruptedAOF = errors.New("corrupted aof file")

// LegacyAOFReader represents the older AOF file reader.
type LegacyAOFReader struct {
	r     io.Reader // reader
	rerr  error     // read error
	chunk []byte    // chunk buffer
	buf   []byte    // main buffer
	l     int       // length of valid data in buffer
	p     int       // pointer
}

// ReadCommand reads an old command.
func (rd *LegacyAOFReader) ReadCommand() ([]byte, error) {
	if rd.l >= 4 {
		sz1 := int(binary.LittleEndian.Uint32(rd.buf[rd.p:]))
		if rd.l >= sz1+9 {
			// we have enough data for a record
			sz2 := int(binary.LittleEndian.Uint32(rd.buf[rd.p+4+sz1:]))
			if sz2 != sz1 || rd.buf[rd.p+4+sz1+4] != 0 {
				return nil, errCorruptedAOF
			}
			buf := rd.buf[rd.p+4 : rd.p+4+sz1]
			rd.p += sz1 + 9
			rd.l -= sz1 + 9
			return buf, nil
		}
	}
	// need more data
	if rd.rerr != nil {
		if rd.rerr == io.EOF {
			rd.rerr = nil // we want to return EOF, but we want to be able to try again
			if rd.l != 0 {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, io.EOF
		}
		return nil, rd.rerr
	}
	if rd.p != 0 {
		// move p to the beginning
		copy(rd.buf, rd.buf[rd.p:rd.p+rd.l])
		rd.p = 0
	}
	var n int
	n, rd.rerr = rd.r.Read(rd.chunk)
	if n > 0 {
		cbuf := rd.chunk[:n]
		if len(rd.buf)-rd.l < n {
			if len(rd.buf) == 0 {
				rd.buf = make([]byte, len(cbuf))
				copy(rd.buf, cbuf)
			} else {
				copy(rd.buf[rd.l:], cbuf[:len(rd.buf)-rd.l])
				rd.buf = append(rd.buf, cbuf[len(rd.buf)-rd.l:]...)
			}
		} else {
			copy(rd.buf[rd.l:], cbuf)
		}
		rd.l += n
	}
	return rd.ReadCommand()
}

// NewLegacyAOFReader creates a new LegacyAOFReader.
func NewLegacyAOFReader(r io.Reader) *LegacyAOFReader {
	rd := &LegacyAOFReader{r: r, chunk: make([]byte, 0xFFFF)}
	return rd
}

func (s *Server) migrateAOF() error {
	_, err := os.Stat(path.Join(s.dir, "appendonly.aof"))
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	_, err = os.Stat(path.Join(s.dir, "aof"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	log.Warn("Migrating aof to new format")
	newf, err := os.Create(path.Join(s.dir, "migrate.aof"))
	if err != nil {
		return err
	}
	defer newf.Close()

	oldf, err := os.Open(path.Join(s.dir, "aof"))
	if err != nil {
		return err
	}
	defer oldf.Close()
	start := time.Now()
	count := 0
	wr := bufio.NewWriter(newf)
	rd := NewLegacyAOFReader(oldf)
	for {
		cmdb, err := rd.ReadCommand()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line := string(cmdb)
		var tok string
		values := make([]resp.Value, 0, 64)
		for line != "" {
			line, tok = token(line)
			if len(tok) > 0 && tok[0] == '{' {
				if line != "" {
					tok = tok + " " + line
					line = ""
				}
			}
			values = append(values, resp.StringValue(tok))
		}
		data, err := resp.ArrayValue(values).MarshalRESP()
		if err != nil {
			return err
		}
		if _, err := wr.Write(data); err != nil {
			return err
		}
		if wr.Buffered() > 1024*1024 {
			if err := wr.Flush(); err != nil {
				return err
			}
		}
		count++
	}
	if err := wr.Flush(); err != nil {
		return err
	}
	oldf.Close()
	newf.Close()
	log.Debugf("%d items: %.0f/sec", count, float64(count)/(float64(time.Since(start))/float64(time.Second)))
	return os.Rename(path.Join(s.dir, "migrate.aof"), path.Join(s.dir, "appendonly.aof"))
}
