package directions

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
	"os"
	"testing"

	"github.com/bhojpur/space/pkg/client/base"
)

func TestDirections(t *testing.T) {

	token := os.Getenv("BHOJPUR_SPACE_MAPS_TOKEN")
	if token == "" {
		t.Error("Maps API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)

	Directions := NewDirections(b)

	t.Run("Can Lookup Directions", func(t *testing.T) {
		var opts RequestOpts

		locs := []base.Location{{37.78, -122.42}, {38.91, -77.03}}

		res, err := Directions.GetDirections(locs, RoutingCycling, &opts)
		if err != nil {
			t.Error(err)
		}

		if Codes(res.Code) != CodeOK {
			t.Errorf("Invalid response code: %s", res.Code)
		}

	})

}
