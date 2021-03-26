// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

import "sync"

var (
	formatMutex   sync.RWMutex
	formatObjects map[string]formatObject
)

func tryParseFormat(source string) formatObject {

}
