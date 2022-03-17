package assert

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
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

var stderr = io.Writer(os.Stderr)
var exit = os.Exit

// Assert aborts the program if assertion is false.
func Assert(expression bool) {
	if !expression {
		atext := "at "
		ftext := ""
		_, file, ln, ok := runtime.Caller(1)
		if ok {
			data, err := os.ReadFile(file)
			if err == nil {
				lines := strings.Split(string(data), "\n")
				if ln-1 < len(lines) {
					if strings.Contains(lines[ln-1], "Assert") {
						atext = strings.Split(lines[ln-1], "Assert")[1]
						if len(atext) > 0 && atext[0] == '(' {
							// create a clean, single line expression
							var cond string
							body := []byte(
								atext + strings.Join(lines[ln:], "\n"),
							)
							depth := 1
							for i := 1; i < len(body); i++ {
								if body[i] == '\n' {
									continue
								} else if body[i] == '\t' {
									body[i] = ' '
								}
								if body[i] == ' ' && body[i-1] == ' ' {
									continue
								}
								cond += string(body[i])
								if body[i] == '"' {
									i++
									for ; i < len(body); i++ {
										cond += string(body[i])
										if body[i] == '"' {
											if body[i-1] == '\\' {
												continue
											}
											break
										}
									}
								} else if body[i] == '(' {
									depth++
								} else if body[i] == ')' {
									depth--
									if depth == 0 {
										cond = cond[:len(cond)-1]
										break
									}
								}
							}
							if cond != "" {
								cond = strings.TrimSpace(cond)
								cond = strings.TrimSuffix(cond, ",")
								atext = "(" + cond + "), "
							}
						}
					}
				}

				fpcs := make([]uintptr, 1)
				// get the function name
				n := runtime.Callers(2, fpcs)
				if n > 0 {
					caller := runtime.FuncForPC(fpcs[0] - 1)
					if caller != nil {
						parts := strings.Split(caller.Name(), "/")
						ftext = parts[len(parts)-1]
						parts = strings.Split(ftext, ".")
						if len(parts) > 1 {
							ftext = strings.Join(parts[1:], ".")
						}
						ftext = strings.Replace(ftext, "(", "", -1)
						ftext = strings.Replace(ftext, ")", "", -1)
						ftext = strings.TrimPrefix(ftext, "*")
						if strings.Contains(ftext, "..") {
							ftext = "[anonymous]"
						}
						ftext = "function " + ftext + ", "
					}
				}
			}
		}
		// make file name relative, if possible
		if path.IsAbs(file) {
			dir, err := os.Getwd()
			if err == nil && strings.HasPrefix(file, dir) {
				tfile := file[len(dir):]
				if len(tfile) > 0 && (tfile[0] == '/' || tfile[0] == '\\') {
					file = tfile[1:]
				}
			}
		}
		fmt.Fprintf(stderr, "Assertion failed: %s%sfile %s, line %d.\n",
			atext, ftext, file, ln)
		exit(6) // SIGABRT
	}
}
