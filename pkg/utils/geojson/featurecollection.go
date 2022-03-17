package geojson

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

	"github.com/bhojpur/space/pkg/utils/gjson"
)

// FeatureCollection ...
type FeatureCollection struct{ collection }

// NewFeatureCollection ...
func NewFeatureCollection(features []Object) *FeatureCollection {
	g := new(FeatureCollection)
	g.children = features
	g.parseInitRectIndex(DefaultParseOptions)
	return g
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g *FeatureCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"FeatureCollection","features":[`...)
	for i := 0; i < len(g.children); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.children[i].AppendJSON(dst)
	}
	dst = append(dst, ']')
	if g.extra != nil {
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	strings.Index("", " ")
	return dst
}

// String ...
func (g *FeatureCollection) String() string {
	return string(g.AppendJSON(nil))
}

// JSON ...
func (g *FeatureCollection) JSON() string {
	return string(g.AppendJSON(nil))
}

// MarshalJSON ...
func (g *FeatureCollection) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func parseJSONFeatureCollection(
	keys *parseKeys, opts *ParseOptions,
) (Object, error) {
	var g FeatureCollection
	if !keys.rFeatures.Exists() {
		return nil, errFeaturesMissing
	}
	if !keys.rFeatures.IsArray() {
		return nil, errFeaturesInvalid
	}
	var err error
	keys.rFeatures.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Parse(value.Raw, opts)
		if err != nil {
			return false
		}
		g.children = append(g.children, f)
		return true
	})
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	g.parseInitRectIndex(opts)
	return &g, nil
}
