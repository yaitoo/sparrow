// +build testable

// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlySimpleVerbs(t *testing.T) {
	verbs := []string{"%b", "%c", "%d", "%e", "%E", "%f", "%F", "%g", "%G", "%o", "%O", "%p", "%q", "%t", "%T", "%U", "%v", "%x", "%X"}

	args := []interface{}{1}

	for _, f := range verbs {
		assert.Equal(t, fmt.Sprintf(f, args...), Sprintf(f, args...), "Simple verb %s must be same", f)
	}

}

func TestOnlyString(t *testing.T) {
	f := "only string without format"

	args := []interface{}{1}

	assert.Equal(t, fmt.Sprintf(f, args...), Sprintf(f, args...), "Only string %s must be same", f)
}

func TestLiteralPercentSign(t *testing.T) {
	f := "only %% string %v"

	args := []interface{}{1}

	assert.Equal(t, fmt.Sprintf(f, args...), Sprintf(f, args...), "LiteralPercentSign %s must be same", f)
}

func TestWithoutVerb(t *testing.T) {
	f := "only string"

	args := []interface{}{1}

	assert.Equal(t, fmt.Sprintf(f, args...), Sprintf(f, args...), "WithoutVerb %s must be same", f)
}

func TestMalformedVerb(t *testing.T) {
	f := "only %y string "

	args := []interface{}{1}

	assert.Equal(t, fmt.Sprintf(f, args...), Sprintf(f, args...), "Malformed verb %s must be same", f)
}

func TestSimpleFormat(t *testing.T) {
	f := "%v"

	args := []interface{}{1}

	assert.Equal(t, fmt.Sprintf(f, args...), Sprintf(f, args...), "Full format verb should be same")

	f2 := "a %v"

	assert.Equal(t, fmt.Sprintf(f2, args...), Sprintf(f2, args...), "String and Format verb should be same")

	f3 := "%v"
	args2 := []interface{}{1, 2}

	assert.Equal(t, fmt.Sprintf(f3, args2...), Sprintf(f3, args2...), "Missing verb should be same")

}
