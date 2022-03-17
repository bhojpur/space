package mapmatching

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/space/pkg/client/base"
)

func TestMapMatching(t *testing.T) {

	token := os.Getenv("BHOJPUR_SPACE_MAPS_TOKEN")
	if token == "" {
		t.Error("Maps API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)

	MapMatching := NewMapMaptching(b)

	timeStamps := []int64{1492878132, 1492878142, 1492878152, 1492878172, 1492878182, 1492878192, 1492878202, 1492878302}
	radiusList := []int{9, 6, 8, 11, 8, 4, 8, 8}

	locs := []base.Location{{37.75319556403746, -122.44254112243651}, {37.75373846204306, -122.44238018989562},
		{37.754111702111146, -122.44199395179749}, {37.75473941979767, -122.44177401065825},
		{37.755570713402115, -122.4412429332733}, {37.756401997666046, -122.44113564491273},
		{37.75677098309616, -122.44228899478911}, {37.756949113334784, -122.4424821138382}}

	t.Run("Map matching supports Polyline", func(t *testing.T) {

		var opts RequestOpts
		opts.SetGeometries(GeometryPolyline)
		opts.SetOverview(OverviewFull)
		opts.SetTimestamps(timeStamps)
		opts.SetSteps(false)
		opts.SetAnnotations([]AnnotationType{AnnotationDistance, AnnotationSpeed})
		opts.SetRadiuses(radiusList)

		res, err := MapMatching.GetMatching(locs, RoutingCycling, &opts)
		assert.Nil(t, err)

		assert.EqualValues(t, Codes(res.Code), CodeOK)

		_, err = res.Matchings[0].GetGeometryPolyline()
		assert.Nil(t, err)

		_, err = res.Matchings[0].GetGeometryGeojson()
		assert.NotNil(t, err)
	})

	t.Run("Map matching supports GeometryGeojson", func(t *testing.T) {

		var opts RequestOpts
		opts.SetGeometries(GeometryGeojson)
		opts.SetOverview(OverviewFull)
		opts.SetTimestamps(timeStamps)
		opts.SetSteps(false)
		opts.SetAnnotations([]AnnotationType{AnnotationDistance, AnnotationSpeed})
		opts.SetRadiuses(radiusList)

		res, err := MapMatching.GetMatching(locs, RoutingCycling, &opts)
		assert.Nil(t, err)

		assert.EqualValues(t, Codes(res.Code), CodeOK)

		_, err = res.Matchings[0].GetGeometryGeojson()
		assert.Nil(t, err)

		_, err = res.Matchings[0].GetGeometryPolyline()
		assert.NotNil(t, err)
	})
}
