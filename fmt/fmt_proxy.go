// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

import (
	"fmt"
	"io"
)

//proxy funcions without any format verb of fmt package

//Errorf a wrapper of fmt.Errorf
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

//Fprint a wrapper of fmt.Fprint
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, a...)
}

//Fprintln a wrapper of fmt.Fprintln
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, a...)
}

//Fscan a wrapper of fmt.Fscan
func Fscan(r io.Reader, a ...interface{}) (n int, err error) {
	return fmt.Fscan(r, a...)
}

//Fscanf a wrapper of fmt.Fscanf
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error) {
	return fmt.Fscanf(r, format, a...)
}

//Fscanln a wrapper of fmt.Fscanln
func Fscanln(r io.Reader, a ...interface{}) (n int, err error) {
	return fmt.Fscanln(r, a...)
}

//Print a wrapper of fmt.Print
func Print(a ...interface{}) (n int, err error) {
	return fmt.Print(a...)
}

//Println a wrapper of fmt.Println
func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

//Scan a wrapper of fmt.Scan
func Scan(a ...interface{}) (n int, err error) {
	return fmt.Scan(a...)
}

//Scanf a wrapper of fmt.Scanf
func Scanf(format string, a ...interface{}) (n int, err error) {
	return fmt.Scanf(format, a...)
}

//Scanln a wrapper of fmt.Scanln
func Scanln(a ...interface{}) (n int, err error) {
	return fmt.Scanln(a...)
}

//Sprint a wrapper of fmt.Sprint
func Sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

//Sprintln a wrapper of fmt.Sprintln
func Sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

//Sscan a wrapper of fmt.Sscan
func Sscan(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscan(str, a...)
}

//Sscanf a wrapper of fmt.Sscanf
func Sscanf(str string, format string, a ...interface{}) (n int, err error) {
	return fmt.Sscanf(str, format, a...)
}

//Sscanln a wrapper of fmt.Sscanln
func Sscanln(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscanln(str, a...)
}
