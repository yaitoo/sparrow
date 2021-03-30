// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"os"
	"sync"
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
		content: content,
		modTime: now,
	}

	c := Open(context.TODO(), "TestWithReader", WithReader(func(ctx context.Context) Reader {
		return r
	}))

	if c.Content != content {
		t.Errorf("Bytes got: %s , want: %s", c.Content, content)
	}

	if c.modTime != now {
		t.Errorf("ModTime got: %v , want: %v", c.modTime, now)
	}
}

func TestOnChanged(t *testing.T) {
	content := "TestOnChanged"
	now := time.Now().UnixNano()

	r := &reader{
		content: content,
		modTime: now,
	}

	c := Open(context.TODO(), "TestOnChanged", WithReader(func(ctx context.Context) Reader {
		return r
	}))

	wantedContent := "Updated:TestOnChanged"
	wantedModTime := time.Now().UnixNano()

	firedHanlder := make(chan bool)

	c.OnChanged(func(c *Config) {
		if wantedContent != c.Content {
			t.Errorf("Bytes got: %s , want: %s", c.Content, wantedContent)
		}

		if c.modTime != wantedModTime {
			t.Errorf("ModTime got: %v , want: %v", c.modTime, wantedModTime)
		}

		firedHanlder <- true
	})

	r.Lock()
	r.content = wantedContent
	r.modTime = wantedModTime
	r.Unlock()

	select {
	case <-time.After(2 * time.Second):
		t.Error("handler is timeout to fired")
	case fired := <-firedHanlder:
		if !fired {
			t.Error("handler is not fired")
		}
	}

}

type reader struct {
	sync.RWMutex
	content string
	modTime int64
}

func (r *reader) Read(ctx context.Context) (string, error) {
	r.RLock()
	defer r.RUnlock()
	return r.content, nil
}

func (r *reader) ModTime(ctx context.Context) (int64, error) {
	r.RLock()
	defer r.RUnlock()
	return r.modTime, nil
}
