// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

//FormatVerb format verb for printing argument
type FormatVerb int

const (
	//FormatVerb_v %v
	FormatVerb_v FormatVerb = iota
	//FormatVerb_T %T
	FormatVerb_T
	//FormatVerb_t %t
	FormatVerb_t
	//FormatVerb_b %b
	FormatVerb_b
	//FormatVerb_c %c
	FormatVerb_c
	//FormatVerb_d %d
	FormatVerb_d
	//FormatVerb_o %o
	FormatVerb_o
	//FormatVerb_O %O
	FormatVerb_O
	//FormatVerb_q %q
	FormatVerb_q
	//FormatVerb_x %x
	FormatVerb_x
	//FormatVerb_X %X
	FormatVerb_X
	//FormatVerb_U %U
	FormatVerb_U
	//FormatVerb_e %e
	FormatVerb_e
	//FormatVerb_E %E
	FormatVerb_E
	//FormatVerb_f %f
	FormatVerb_f
	//FormatVerb_F %F
	FormatVerb_F
	//FormatVerb_g %g
	FormatVerb_g
	//FormatVerb_G %G
	FormatVerb_G
	//FormatVerb_p %p
	FormatVerb_p
)

func (f FormatVerb) String() string {
	return [...]string{"%v", "%T", "%t", "%b", "%c", "%d", "%o", "%O", "%q", "%x", "%X", "%U", "%e", "%E", "%f", "%F", "%g", "%G", "%p"}[f]
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
	stringParts []string
	formatParts []formatIndex
}

//formatIndex the metedata of format part, include index, source, flags and print function
type formatIndex struct {
	index  int
	source string
	flags  FormatFlags
	print  PrintFunc
}
