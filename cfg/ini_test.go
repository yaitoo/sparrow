// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cfg

import (
	"testing"
)

func TestSections(t *testing.T) {

	i := Inifile{}
	i.TryParse(`
	#dev mysql server
	[mysql]
	host=127.0.0.1:3306
	#dev redis server
	[redis]
	host=127.0.0.1:6379
	`)

	wantedMySQLHost := "127.0.0.1:3306"
	acutalMySQLHost := i.Section("mysql").Value("host", "")

	if acutalMySQLHost != wantedMySQLHost {
		t.Errorf("[mysql]host got: %s, want: %s", acutalMySQLHost, wantedMySQLHost)
	}

	wantedRedisHost := "127.0.0.1:6379"
	acutalRedisHost := i.Section("redis").Value("host", "")

	if acutalMySQLHost != wantedMySQLHost {
		t.Errorf("[redis]host got: %s, want: %s", acutalRedisHost, wantedRedisHost)
	}
}

func TestNormalize(t *testing.T) {

	i := Inifile{}
	i.TryParse(`
	#dev mysql server
[ MySQL ]
Host =127.0.0.1:3306
	#dev redis server
 [ Redis]
	host =127.0.0.1:6379
	`)

	wantedMySQLHost := "127.0.0.1:3306"
	acutalMySQLHost := i.Section("mysql").Value("host", "")

	if acutalMySQLHost != wantedMySQLHost {
		t.Errorf("[mysql]host got: %s, want: %s", acutalMySQLHost, wantedMySQLHost)
	}

	wantedRedisHost := "127.0.0.1:6379"
	acutalRedisHost := i.Section("redis").Value("host", "")

	if acutalMySQLHost != wantedMySQLHost {
		t.Errorf("[redis]host got: %s, want: %s", acutalRedisHost, wantedRedisHost)
	}
}

func TestValues(t *testing.T) {

	i := Inifile{}
	i.TryParse(`
[values]
string=s
int=3
int32=32
int64=64
float32= 3.2  
float64= 6.4  
bool_0=0
bool_1=1
bool_on=on
bool_off=off
bool_true=true
bool_false=false`)

	s := i.Section("values")
	if s.Value("string", "") != "s" {
		t.Errorf("Value got: %s, want: %s", s.Value("string", ""), "s")
	}

	if s.ValueInt("int", 0) != 3 {
		t.Errorf("ValueInt got: %v, want: %v", s.ValueInt("int", 0), 3)
	}

	if s.ValueInt32("int32", 0) != 32 {
		t.Errorf("ValueInt32 got: %v, want: %v", s.ValueInt32("int32", 0), 32)
	}

	if s.ValueInt64("int64", 0) != 64 {
		t.Errorf("ValueInt32 got: %v, want: %v", s.ValueInt64("int64", 0), 64)
	}

	if s.ValueFloat32("float32", 0) != 3.2 {
		t.Errorf("ValueFloat32 got: %v, want: %v", s.ValueFloat32("float32", 0), 3.2)
	}

	if s.ValueFloat64("float64", 0) != 6.4 {
		t.Errorf("ValueFloat64 got: %v, want: %v", s.ValueFloat64("float64", 0), 6.4)
	}

	if s.ValueBool("bool_0", true) != false {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_0", true), false)
	}

	if s.ValueBool("bool_1", false) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_1", false), true)
	}

	if s.ValueBool("bool_off", true) != false {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_off", true), false)
	}

	if s.ValueBool("bool_on", false) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_on", false), true)
	}

	if s.ValueBool("bool_false", true) != false {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_false", true), false)
	}

	if s.ValueBool("bool_true", false) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_true", false), true)
	}

}

func TestNormailizeValues(t *testing.T) {

	i := Inifile{}
	i.TryParse(`
[bool]
bool_1= 1
bool_on= On
bool_off= Off
bool_true= True
bool_false= False
bool_invalid= `)

	s := i.Section("bool")

	if s.ValueBool("bool_1", false) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_1", false), true)
	}

	if s.ValueBool("bool_off", true) != false {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_off", true), false)
	}

	if s.ValueBool("bool_on", false) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_on", false), true)
	}

	if s.ValueBool("bool_false", true) != false {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_false", true), false)
	}

	if s.ValueBool("bool_true", false) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_true", false), true)
	}

	if s.ValueBool("bool_invalid", true) != true {
		t.Errorf("ValueBool got: %v, want: %v", s.ValueBool("bool_invalid", true), true)
	}

}
