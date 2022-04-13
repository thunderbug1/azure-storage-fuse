/*
    _____           _____   _____   ____          ______  _____  ------
   |     |  |      |     | |     | |     |     | |       |            |
   |     |  |      |     | |     | |     |     | |       |            |
   | --- |  |      |     | |-----| |---- |     | |-----| |-----  ------
   |     |  |      |     | |     | |     |     |       | |       |
   | ____|  |_____ | ____| | ____| |     |_____|  _____| |_____  |_____


   Licensed under the MIT License <http://opensource.org/licenses/MIT>.

   Copyright © 2020-2022 Microsoft Corporation. All rights reserved.
   Author : <blobfusedev@microsoft.com>

   Permission is hereby granted, free of charge, to any person obtaining a copy
   of this software and associated documentation files (the "Software"), to deal
   in the Software without restriction, including without limitation the rights
   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
   copies of the Software, and to permit persons to whom the Software is
   furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all
   copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
   SOFTWARE
*/

package attr_cache

import (
	"blobfuse2/common"
	"blobfuse2/internal"
	"os"
	"time"
)

// Flags represented in BitMap for various flags in the attr cache item
const (
	AttrFlagUnknown uint16 = iota
	AttrFlagExists
	AttrFlagValid
)

// attrCacheItem : Structure of each item in attr cache
type attrCacheItem struct {
	attr     *internal.ObjAttr
	cachedAt time.Time
	attrFlag common.BitMap16
}

func newAttrCacheItem(attr *internal.ObjAttr, exists bool, cachedAt time.Time) *attrCacheItem {
	item := &attrCacheItem{
		attr:     attr,
		attrFlag: 0,
		cachedAt: cachedAt,
	}

	item.attrFlag.Set(AttrFlagValid)
	if exists {
		item.attrFlag.Set(AttrFlagExists)
	}

	return item
}

func (item *attrCacheItem) valid() bool {
	return item.attrFlag.IsSet(AttrFlagValid)
}

func (item *attrCacheItem) exists() bool {
	return item.attrFlag.IsSet(AttrFlagExists)
}

func (item *attrCacheItem) markDeleted(deletedTime time.Time) {
	item.attrFlag.Clear(AttrFlagExists)
	item.attrFlag.Set(AttrFlagValid)
	item.cachedAt = deletedTime
	item.attr = &internal.ObjAttr{}
}

func (item *attrCacheItem) invalidate() {
	item.attrFlag.Clear(AttrFlagValid)
	item.attr = &internal.ObjAttr{}
}

func (item *attrCacheItem) getAttr() *internal.ObjAttr {
	return item.attr
}

func (item *attrCacheItem) isDeleted() bool {
	return !item.exists()
}

func (item *attrCacheItem) setSize(size int64) {
	item.attr.Mtime = time.Now()
	item.attr.Size = size
	item.cachedAt = time.Now()
}

func (item *attrCacheItem) setMode(mode os.FileMode) {
	item.attr.Mode = mode
	item.attr.Ctime = time.Now()
	item.cachedAt = time.Now()
}

func (item *attrCacheItem) setMetadata(key string, value string) {
	item.attr.Metadata[key] = value
	item.attr.Ctime = time.Now()
	item.cachedAt = time.Now()
}

func (item *attrCacheItem) removeMetadata(key string) {
	delete(item.attr.Metadata, key)
	item.attr.Ctime = time.Now()
	item.cachedAt = time.Now()
}
