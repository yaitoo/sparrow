package validator

import (
	"regexp"
)

/*
	Required
	Min(min int)
	Max(max int)
	Range(min, max int)
	MinSize(min int)
	MaxSize(max int)
	Length(length int)
	Alpha
	Numeric
	AlphaNumeric
	Match(pattern string)
	AlphaDash
	Email
	IP
	Base64
	Mobile
	Tel
	Phone
	ZipCode
*/

const (
	wordsize = 32 << (^uint(0) >> 32 & 1)
)

var (
	alphaDashPattern = regexp.MustCompile(`[^\d\w-_]`)
	emailPattern     = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)
	ipPattern        = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	base64Pattern    = regexp.MustCompile(`^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`)
	// just for chinese mobile phone number
	mobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][06789]|[4][579]))\d{8}$`)
	// just for chinese telephone number
	telPattern = regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`)
	// just for chinese zipcode
	zipCodePattern = regexp.MustCompile(`^[1-9]\d{5}$`)
)
