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

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	// BaseURL Maps API base URL
	BaseURL = "https://maps.bhojpur.net"

	statusRateLimitExceeded = 429
)

// Base Maps API base
type Base struct {
	token string
	debug bool
}

// NewBase Create a new API base instance
func NewBase(token string) *Base {
	b := &Base{}

	b.token = token

	return b
}

// SetDebug enables debug output for API calls
func (b *Base) SetDebug(debug bool) {
	b.debug = true
}

type MapEngineApiMessage struct {
	Message string
}

// QueryRequest make a get with the provided query string and return the response if successful
func (b *Base) QueryRequest(query string, v *url.Values) (*http.Response, error) {
	// Add token to args
	v.Set("access_token", b.token)

	// Generate URL
	url := fmt.Sprintf("%s/%s", BaseURL, query)

	if b.debug {
		fmt.Printf("URL: %s\n", url)
	}

	// Create request object
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = v.Encode()

	// Create client instance
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if b.debug {
		data, _ := httputil.DumpRequest(request, true)
		fmt.Printf("Request: %s", string(data))
		data, _ = httputil.DumpResponse(resp, false)
		fmt.Printf("Response: %s", string(data))
	}

	if resp.StatusCode == statusRateLimitExceeded {
		return nil, ErrorAPILimitExceeded
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrorAPIUnauthorized
	}

	return resp, nil
}

// QueryBase Query the Maps API and fill the provided instance with the returned JSON
// TODO: Rename this
func (b *Base) QueryBase(query string, v *url.Values, inst interface{}) error {

	resp, err := b.QueryRequest(query, v)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&inst)
	if err != nil {
		return err
	}

	return nil
}

// Query the Maps API
// TODO: Depreciate this
func (b *Base) Query(api, version, mode, query string, v *url.Values, inst interface{}) error {

	// Generate URL
	queryString := fmt.Sprintf("%s/%s/%s/%s", api, version, mode, query)

	return b.QueryBase(queryString, v, inst)
}
