// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

import (
	"bytes"
	"reflect"
	"strconv"
	"sync"
)

var (
	formatMutex   sync.RWMutex
	formatObjects = make(map[string]formatObject)
)

//FormatVerb format verb for printing argument
type FormatVerb int

const (

	//FormatVerbLowerB %b
	FormatVerbLowerB FormatVerb = iota

	//FormatVerbLowerC %c
	FormatVerbLowerC

	//FormatVerbLowerD %d
	FormatVerbLowerD

	//FormatVerbLowerE %e
	FormatVerbLowerE
	//FormatVerbUpperE %E
	FormatVerbUpperE

	//FormatVerbLowerF %f
	FormatVerbLowerF
	//FormatVerbUpperF %F
	FormatVerbUpperF

	//FormatVerbLowerG %g
	FormatVerbLowerG
	//FormatVerbUpperG %G
	FormatVerbUpperG

	//FormatVerbLowerO %o
	FormatVerbLowerO
	//FormatVerbUpperO %O
	FormatVerbUpperO

	//FormatVerbLowerP %p
	FormatVerbLowerP

	//FormatVerbLowerQ %q
	FormatVerbLowerQ

	//FormatVerbLowerT %t
	FormatVerbLowerT
	//FormatVerbUpperT %T
	FormatVerbUpperT

	//FormatVerbUpperU %U
	FormatVerbUpperU

	//FormatVerbLowerV %v
	FormatVerbLowerV

	//FormatVerbLowerX %x
	FormatVerbLowerX
	//FormatVerbUpperX %X
	FormatVerbUpperX
)

var verbs = []string{"%b", "%c", "%d", "%e", "%E", "%f", "%F", "%g", "%G", "%o", "%O", "%p", "%q", "%t", "%T", "%U", "%v", "%x", "%X"}

func (f FormatVerb) String() string {
	return verbs[f]
}

//FormatFlags https://golang.org/pkg/fmt/#hdr-Printing
type FormatFlags struct {
	//Width width. set -1 if it is unset.
	Width int
	//Precision Precision. set -1 if it is unset.
	Precision int
	//Plus always print a sign for numeric values; guarantee ASCII-only output for %q (%+q)
	Plus bool
	//Minus pad with spaces on the right rather than the left (left-justify the field)
	Minus bool
	//Sharp alternate format: add leading 0b for binary (%#b), 0 for octal (%#o),
	// 0x or 0X for hex (%#x or %#X); suppress 0x for %p (%#p);
	// for %q, print a raw (backquoted) string if strconv.CanBackquote
	// returns true;
	// always print a decimal point for %e, %E, %f, %F, %g and %G;
	// do not remove trailing zeros for %g and %G;
	// write e.g. U+0078 'x' if the character is printable for %U (%#U).
	Sharp bool
	//Space (space) leave a space for elided sign in numbers (% d);
	//put spaces between bytes printing strings or slices in hex (% x, % X)
	Space bool
	//Zero pad with leading zeros rather than spaces;
	//for numbers, this moves the padding after the sign
	Zero bool
}

//formatObject parse format, and cache string part and format part in a formatObject
type formatObject struct {
	items       []bool
	stringParts [][]byte
	formatParts []formatIndex
}

func (fo formatObject) Printf(args ...interface{}) string {

	buf := bytes.Buffer{}

	stringIndex := 0
	formatIndex := 0
	argsIndex := 0

	argsLen := len(args)

	n := len(fo.items)
	for i := 0; i < n; i++ {
		if fo.items[i] { //format part

			fi := fo.formatParts[formatIndex]

			if argsIndex < argsLen {
				buf.WriteString(fi.print(fi.source, fi.verb, fi.flags, args[argsIndex]))
			} else {
				buf.WriteString(fi.print(fi.source, fi.verb, fi.flags, nil))
			}

			argsIndex++
			formatIndex++

		} else { //string part
			buf.Write(fo.stringParts[stringIndex])
			stringIndex++
		}
	}

	return buf.String()

}

//formatIndex the metedata of format part, include index, source, flags and print function
type formatIndex struct {
	source string
	verb   FormatVerb
	flags  FormatFlags
	print  PrintfFunc
}

func parseFormatObject(source string, args ...interface{}) formatObject {
	formatMutex.RLock()

	key := strconv.Itoa(len(args)) + ":" + source

	f, ok := formatObjects[key]
	formatMutex.RUnlock()
	if ok {
		return f
	}

	f = formatObject{
		items:       make([]bool, 0, 10),
		stringParts: make([][]byte, 0, 10),
		formatParts: make([]formatIndex, 0, 10),
	}
	end := len(source)
	argNum := 0
	argLen := len(args)

	for i := 0; i < end; {

		buf := bytes.Buffer{}
		lasti := i

		for i < end {

			if source[i] != '%' {
				i++
			} else {
				if source[i+1] == '%' {
					i = i + 2
				} else {
					break
				}
			}
		}

		if i > lasti {
			buf.WriteString(source[lasti:i])
		}

		if i >= end {
			// done processing format string

			f.items = append(f.items, false)
			f.stringParts = append(f.stringParts, buf.Bytes())

			break
		}

		fi := formatIndex{}
		for i < end {
			i++
			c := source[i]
			if 'A' <= c && c <= 'z' {
				fi.source = source[lasti : i+1]

				fi.verb, fi.flags = parseFormatVerb(source)

				if argNum < argLen {
					fi.print = getPrintfFunc(fi.verb, reflect.TypeOf(args[argNum]))
				} else {
					fi.print = defaultPrintfFunc
				}
				argNum++

				f.items = append(f.items, true)
				f.formatParts = append(f.formatParts, fi)

				i++
				lasti = i
				break
			}

		}

		//format % is missing verb at end of string
		if i > lasti {
			fi.source = source[lasti:i]

			fi.print = defaultPrintfFunc

			f.items = append(f.items, true)
			f.formatParts = append(f.formatParts, fi)

			argNum++
		}

	}

	for ; argNum < argLen; argNum++ {
		fi := formatIndex{}
		fi.source = ""
		fi.print = defaultPrintfFunc

		f.items = append(f.items, true)
		f.formatParts = append(f.formatParts, fi)
	}

	formatMutex.Lock()
	formatObjects[key] = f
	formatMutex.Unlock()
	return f
}

func parseFormatVerb(source string) (FormatVerb, FormatFlags) {
	// it := formatIndex{}
	// c := source[i]
	// switch c {
	// case '#':
	// 	it.flags.Sharp = true
	// case '0':
	// 	it.flags.Zero = !it.flags.Minus // Only allow zero padding to the left.
	// case '+':
	// 	it.flags.Plus = true
	// case '-':
	// 	it.flags.Minus = true
	// 	it.flags.Zero = false // Do not pad with zeros to the right.
	// case ' ':
	// 	it.flags.Space = true
	// default:
	// 	// Fast path for common case of ascii lower case simple verbs
	// 	// without precision or width or argument indices.
	// 	if 'a' <= c && c <= 'z' && argNum < len(a) {
	// 		if c == 'v' {
	// 			// Go syntax
	// 			p.fmt.sharpV = p.fmt.sharp
	// 			p.fmt.sharp = false
	// 			// Struct-field syntax
	// 			p.fmt.plusV = p.fmt.plus
	// 			p.fmt.plus = false
	// 		}
	// 		p.printArg(a[argNum], rune(c))
	// 		argNum++
	// 		i++
	// 		continue formatLoop
	// 	}
	// 	// Format is more complex than simple flags and a verb or is malformed.
	// 	break simpleFormat
	// }
	return -1, FormatFlags{}
}
