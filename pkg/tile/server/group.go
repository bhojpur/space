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
	"github.com/bhojpur/space/pkg/utils/btree"
)

func byGroupHook(va, vb interface{}) bool {
	a, b := va.(*groupItem), vb.(*groupItem)
	if a.hookName < b.hookName {
		return true
	}
	if a.hookName > b.hookName {
		return false
	}
	if a.colKey < b.colKey {
		return true
	}
	if a.colKey > b.colKey {
		return false
	}
	return a.objID < b.objID
}

func byGroupObject(va, vb interface{}) bool {
	a, b := va.(*groupItem), vb.(*groupItem)
	if a.colKey < b.colKey {
		return true
	}
	if a.colKey > b.colKey {
		return false
	}
	if a.objID < b.objID {
		return true
	}
	if a.objID > b.objID {
		return false
	}
	return a.hookName < b.hookName
}

type groupItem struct {
	hookName string
	colKey   string
	objID    string
	groupID  string
}

func newGroupItem(hookName, colKey, objID string) *groupItem {
	groupID := bsonID()
	g := &groupItem{}
	// create a single string allocation
	ustr := hookName + colKey + objID + groupID
	var pos int
	g.hookName = ustr[pos : pos+len(hookName)]
	pos += len(hookName)
	g.colKey = ustr[pos : pos+len(colKey)]
	pos += len(colKey)
	g.objID = ustr[pos : pos+len(objID)]
	pos += len(objID)
	g.groupID = ustr[pos : pos+len(groupID)]
	pos += len(groupID)
	return g
}

func (s *Server) groupConnect(hookName, colKey, objID string) (groupID string) {
	g := newGroupItem(hookName, colKey, objID)
	s.groupHooks.Set(g)
	s.groupObjects.Set(g)
	return g.groupID
}

func (s *Server) groupDisconnect(hookName, colKey, objID string) {
	g := &groupItem{
		hookName: hookName,
		colKey:   colKey,
		objID:    objID,
	}
	s.groupHooks.Delete(g)
	s.groupObjects.Delete(g)
}

func (s *Server) groupGet(hookName, colKey, objID string) (groupID string) {
	v := s.groupHooks.Get(&groupItem{
		hookName: hookName,
		colKey:   colKey,
		objID:    objID,
	})
	if v != nil {
		return v.(*groupItem).groupID
	}
	return ""
}

func deleteGroups(s *Server, groups []*groupItem) {
	var hhint btree.PathHint
	var ohint btree.PathHint
	for _, g := range groups {
		s.groupHooks.DeleteHint(g, &hhint)
		s.groupObjects.DeleteHint(g, &ohint)
	}
}

// groupDisconnectObject disconnects all hooks from provide object
func (s *Server) groupDisconnectObject(colKey, objID string) {
	var groups []*groupItem
	s.groupObjects.Ascend(&groupItem{colKey: colKey, objID: objID},
		func(v interface{}) bool {
			g := v.(*groupItem)
			if g.colKey != colKey || g.objID != objID {
				return false
			}
			groups = append(groups, g)
			return true
		},
	)
	deleteGroups(s, groups)
}

// groupDisconnectCollection disconnects all hooks from objects in provided
// collection.
func (s *Server) groupDisconnectCollection(colKey string) {
	var groups []*groupItem
	s.groupObjects.Ascend(&groupItem{colKey: colKey},
		func(v interface{}) bool {
			g := v.(*groupItem)
			if g.colKey != colKey {
				return false
			}
			groups = append(groups, g)
			return true
		},
	)
	deleteGroups(s, groups)
}

// groupDisconnectHook disconnects all objects from provided hook.
func (s *Server) groupDisconnectHook(hookName string) {
	var groups []*groupItem
	s.groupHooks.Ascend(&groupItem{hookName: hookName},
		func(v interface{}) bool {
			g := v.(*groupItem)
			if g.hookName != hookName {
				return false
			}
			groups = append(groups, g)
			return true
		},
	)
	deleteGroups(s, groups)
}
