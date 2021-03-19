// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//Package cfg implements a simple config loader with hot reload feature.
//Local and remote file are both supported.
package cfg

import (
	"context"
	"sync"
)

var (
	configs = make(map[string]*Config)
	mutex   sync.RWMutex
)

//Handler A hanlder responds to a config changed.
type Handler func(c *Config)

//Option set optional parameter
type Option func(ctx context.Context, c *Config)

//Config Config represents an open config descriptor.
type Config struct {
	sync.RWMutex
	//Name name
	Name string
	//Bytes content
	Bytes []byte

	//modTime modification time
	modTime int64

	//reader a Reader instance
	reader Reader
	//handlers handlers respond to config changed
	hanlders []Handler
}

//OnChanged subscribe config change
func (c *Config) OnChanged(handler Handler) {

	if handler == nil {
		return
	}

	c.Lock()
	defer c.Unlock()

	if c.hanlders == nil {
		c.hanlders = []Handler{handler}
	} else {
		c.hanlders = append(c.hanlders, handler)
	}
}

func (c *Config) fireHandlers() {
	c.RLock()
	defer c.RUnlock()

	for _, handler := range c.hanlders {
		go handler(c)
	}

}

//ToInifile convert config to Inifile
func (c *Config) ToInifile() Inifile {
	i := Inifile{}
	i.TryParse(string(c.Bytes))

	return i
}

//Open open a Config with name
func Open(ctx context.Context, name string, options ...Option) *Config {
	mutex.Lock()
	defer mutex.Unlock()

	c, ok := configs[name]

	if ok {
		return c
	}

	c = &Config{
		Name: name,
	}

	configs[name] = c

	for _, option := range options {
		option(ctx, c)
	}

	if c.reader == nil {
		c.reader = CreateFsReader(ctx, name)
	}

	c.Bytes, _ = c.reader.Read(ctx)
	c.modTime, _ = c.reader.ModTime(ctx)

	once.Do(startWatch)

	return c
}

//WithReader using a custom Reader
func WithReader(create func(ctx context.Context) Reader) Option {
	return func(ctx context.Context, c *Config) {
		c.reader = create(ctx)
	}
}
