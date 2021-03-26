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

//Errorf formats according to a format specifier and returns the string as a value that satisfies error.
//
//If the format specifier includes a %w verb with an error operand, the returned error will implement an Unwrap method returning the operand. It is invalid to include more than one %w verb or to supply it with an operand that does not implement the error interface. The %w verb is otherwise a synonym for %v.
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, a...)
}
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {

}
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {

}
func Fscan(r io.Reader, a ...interface{}) (n int, err error) {

}
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error) {

}
func Fscanln(r io.Reader, a ...interface{}) (n int, err error) {

}
func Print(a ...interface{}) (n int, err error) {

}
func Printf(format string, a ...interface{}) (n int, err error) {

}
func Println(a ...interface{}) (n int, err error) {

}
func Scan(a ...interface{}) (n int, err error) {

}
func Scanf(format string, a ...interface{}) (n int, err error) {

}
func Scanln(a ...interface{}) (n int, err error) {

}
func Sprint(a ...interface{}) string {

}
func Sprintf(format string, a ...interface{}) string {

}
func Sprintln(a ...interface{}) string {

}
func Sscan(str string, a ...interface{}) (n int, err error) {

}
func Sscanf(str string, format string, a ...interface{}) (n int, err error) {

}
func Sscanln(str string, a ...interface{}) (n int, err error) {

}
