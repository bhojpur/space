package buffer

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
	"testing"

	"github.com/bhojpur/space/pkg/utils/geojson"
	"github.com/bhojpur/space/pkg/utils/geojson/geometry"
)

const lineString = `{"type":"LineString","coordinates":[
	[-116.40289306640624,34.125447565116126],
	[-116.36444091796875,34.14818102254435],
	[-116.0980224609375,34.15045403191448],
	[-115.74920654296874,34.127721186043985],
	[-115.54870605468749,34.075412438417395],
	[-115.5267333984375,34.11407854333859],
	[-115.21911621093749,34.048108084909835],
	[-115.25207519531249,33.8339199536547],
	[-115.40588378906249,33.71748624018193]
]}`

var lineInPoints = []geometry.Point{
	{X: -115.64363479614258, Y: 34.108251327293296},
	{X: -115.54355621337892, Y: 34.07199987534163},
	{X: -115.21482467651367, Y: 34.051237154976164},
	{X: -115.4110336303711, Y: 33.715201644740844},
	{X: -116.40701293945311, Y: 34.12345809664606},
}

func TestBufferLineString(t *testing.T) {
	g, err := geojson.Parse(lineString, nil)
	if err != nil {
		t.Fatal(err)
	}
	g2, err := Simple(g, 1000)
	if err != nil {
		t.Fatal(err)
	}
	for _, pt := range lineInPoints {
		ok := g2.Contains(geojson.NewPoint(pt))
		if !ok {
			t.Fatalf("!ok")
		}
	}
}

const polygon = `{"type": "Polygon","coordinates":[
	[
		[116.46881103515624,34.277644878733824],
		[115.87280273437499,34.20953080048952],
		[115.70251464843749,34.397844946449865],
		[115.9881591796875,34.61286625296406],
		[116.46881103515624,34.277644878733824]
	],
	[
		[115.90438842773436,34.38651267795365],
		[116.05270385742188,34.35023911062779],
		[115.99914550781249,34.44655621402982],
		[115.90438842773436,34.38651267795365]
	]
]}`

var polyInPoints = []geometry.Point{
	{X: 115.95837593078612, Y: 34.59887847065301},
	{X: 115.98755836486816, Y: 34.61879975173954},
	{X: 115.98833084106445, Y: 34.59795999847678},
	{X: 116.04536533355714, Y: 34.58082509817638},
	{X: 116.47567749023438, Y: 34.27651009584797},
	{X: 116.42005920410155, Y: 34.32018817684490},
	{X: 116.33216857910156, Y: 34.25948651450623},
	{X: 115.89340209960939, Y: 34.24132422972854},
	{X: 115.95588684082033, Y: 34.42786803680155},
	{X: 115.97236633300783, Y: 34.42107129982385},
	{X: 115.99639892578125, Y: 34.43579686485573},
	{X: 116.04652404785155, Y: 34.35364042469895},
	{X: 115.92155456542967, Y: 34.38877925439021},
	{X: 115.96755981445311, Y: 34.37687904351907},
	{X: 115.88859558105467, Y: 34.42956713470528},
	{X: 115.97511291503906, Y: 34.36327673174518},
	{X: 115.69564819335938, Y: 34.39784494644986},
	{X: 115.87005615234375, Y: 34.20385213966983},
	{X: 115.76980590820312, Y: 34.31678550602221},
}
var polyOutPoints = []geometry.Point{
	{X: 115.68534851074217, Y: 34.40917568058836},
	{X: 115.98953247070312, Y: 34.63038297923298},
	{X: 115.98541259765624, Y: 34.39671178864245},
	{X: 116.31500244140626, Y: 34.22145474280257},
	{X: 115.85426330566406, Y: 34.18510984477340},
}

func TestBufferPolygon(t *testing.T) {
	g, err := geojson.Parse(polygon, nil)
	if err != nil {
		t.Fatal(err)
	}
	g2, err := Simple(g, 1000)
	if err != nil {
		t.Fatal(err)
	}
	for _, pt := range polyInPoints {
		ok := g2.Contains(geojson.NewPoint(pt))
		if !ok {
			t.Fatalf("!ok")
		}
	}
	for _, pt := range polyOutPoints {
		ok := g2.Contains(geojson.NewPoint(pt))
		if ok {
			t.Fatalf("ok")
		}
	}
}
