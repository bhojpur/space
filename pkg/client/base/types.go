package base

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

type Point []float64

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type BoundingBox []float64

type Geometry struct {
	Type        string `json:"type"`
	Coordinates Point  `json:"coordinates"`
}

type Context struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	ShortCode string `json:"short_code"`
	WikiData  string `json:"wikidata"`
}

type Properties struct {
	Category string `json:"category"`
	Tel      string `json:"tel"`
	Wikidata string `json:"wikidata"`
	Landmark bool   `json:"landmark"`
	Maki     string `json:"short_code"`
}

type Feature struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Text       string      `json:"text"`
	PlaceName  string      `json:"place_name"`
	PlaceType  []string    `json:"place_type"`
	Relevance  float64     `json:"relevance"`
	Properties Properties  `json:"properties"`
	BBox       BoundingBox `json:"bbox"`
	Center     Point       `json:"center"`
	Geometry   Geometry    `json:"geometry"`
	Context    []Context   `json:"context"`
}

type FeatureCollection struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
}
