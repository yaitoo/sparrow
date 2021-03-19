// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	ReloadInterval = 1 * time.Second

	code := m.Run()

	os.Exit(code)
}

func TestOpenWithReader(t *testing.T) {
	content := "TestWithReader"
	now := time.Now().UnixNano()

	r := &reader{
		bytes:   []byte(content),
		modTime: now,
	}

	c := Open(context.TODO(), "TestWithReader", WithReader(func(ctx context.Context, c *Config) Reader {
		return r
	}))

	if string(c.Bytes) != content {
		t.Errorf("Bytes got: %s , want: %s", string(c.Bytes), content)
	}

	if c.modTime != now {
		t.Errorf("ModTime got: %v , want: %v", c.modTime, now)
	}
}

func TestOnChanged(t *testing.T) {
	content := "TestOnChanged"
	now := time.Now().UnixNano()

	r := &reader{
		bytes:   []byte(content),
		modTime: now,
	}

	c := Open(context.TODO(), "TestOnChanged", WithReader(func(ctx context.Context, c *Config) Reader {
		return r
	}))

	wantedContent := "Updated:TestOnChanged"
	wantedModTime := time.Now().UnixNano()

	firedHanlder := false

	go c.OnChanged(func(c *Config) {
		if wantedContent != string(c.Bytes) {
			t.Errorf("Bytes got: %s , want: %s", string(c.Bytes), wantedContent)
		}

		if c.modTime != wantedModTime {
			t.Errorf("ModTime got: %v , want: %v", c.modTime, wantedModTime)
		}

		firedHanlder = true
	})

	r.bytes = []byte(wantedContent)
	r.modTime = wantedModTime

	time.Sleep(2 * time.Second)

	if !firedHanlder {
		t.Error("handler is not fired")
	}

}

type reader struct {
	bytes   []byte
	modTime int64
}

func (r *reader) Read(ctx context.Context) ([]byte, error) {
	return r.bytes, nil
}

func (r *reader) ModTime(ctx context.Context) (int64, error) {
	return r.modTime, nil
}
