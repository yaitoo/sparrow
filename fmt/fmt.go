// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//Package fmt implements a simple fmt wrapper with higher performance.
//It will hit performance issue when `fmt.Sprintf`,`fmt.Printf` or `fmt.Fprintf` is caled too many times.Because it works based on `reflect`.
//The package wraps standard `fmt`, and improves performance by caching and reusing reflect result.
package fmt

import (
	"fmt"
	"io"
)

//Fprintf a high-performance Fprintf instead of fmt.Fprintf
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	fo := parseFormatObject(format, a...)
	return fmt.Fprint(w, fo.Printf(a...))
}

//Sprintf a high-performance Sprintf instead of fmt.Sprintf
func Sprintf(format string, a ...interface{}) string {

	fo := parseFormatObject(format, a...)

	return fo.Printf(a...)
}

//Printf a high-performance Printf instead of fmt.Printf
func Printf(format string, a ...interface{}) (n int, err error) {

	fo := parseFormatObject(format, a...)

	return fmt.Print(fo.Printf(a...))
}
