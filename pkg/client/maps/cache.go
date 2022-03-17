package maps

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
	"image"
	"os"
	"strings"
)

// FileCache is a simple file-based caching implementation for map tiles
// This does not implement any mechanisms for deletion / removal, and as such is not suitable for production use
type FileCache struct {
	basePath string
}

// NewFileCache creates a new file cache instance
func NewFileCache(basePath string) (*FileCache, error) {
	fc := &FileCache{basePath}

	err := os.Mkdir(basePath, 0777)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return fc, nil
}

func (fc *FileCache) getName(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool) string {
	dpiString := ""
	if highDPI {
		dpiString = "@2x"
	}
	return fmt.Sprintf("%s-%d-%d-%d%s.%s", mapID, x, y, level, dpiString, format)
}

// Save saves an image to the file cache
func (fc *FileCache) Save(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool, img image.Image) error {
	name := fc.getName(mapID, x, y, level, format, highDPI)
	path := fmt.Sprintf("%s/%s", fc.basePath, name)

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	// TODO: cannot currently save pngraw
	if strings.Contains(string(format), "png") {
		return SaveImagePNG(img, path)
	}

	if strings.Contains(string(format), "jpg") || strings.Contains(string(format), "jpeg") {
		return SaveImageJPG(img, path)
	}

	return fmt.Errorf("Unrecognized file type (%s)", format)
}

// Fetch fetches an image from the file cache if possible
func (fc *FileCache) Fetch(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool) (image.Image, *image.Config, error) {
	name := fc.getName(mapID, x, y, level, format, highDPI)
	path := fmt.Sprintf("%s/%s", fc.basePath, name)

	// TODO: work out how to load raw png files :-/
	if format == MapFormatPngRaw {
		return nil, nil, nil
	}

	if _, err := os.Stat(path); err != nil {
		return nil, nil, nil
	}

	img, cfg, err := LoadImage(path)

	return img, cfg, err
}
