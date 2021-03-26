// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Reader Reader is the interface that wrap the Read and ModTime method of config
type Reader interface {
	//Read reads contents
	Read(ctx context.Context) (string, error)
	//ModTime get latest modification time
	ModTime(ctx context.Context) (int64, error)
}

//FsReader FsReader implements `Reader` inferface for local file system
type FsReader struct {
	fileName string
}

//CreateFsReader create an instance of FsReader
func CreateFsReader(ctx context.Context, name string) Reader {

	f := &FsReader{}

	fileName, err := filepath.Abs(name)

	if err == nil {
		f.fileName = name
	} else {
		f.fileName = fileName
	}

	return f
}

//Read implement `Reader.Read`
func (r *FsReader) Read(ctx context.Context) (string, error) {
	buf, err := ioutil.ReadFile(r.fileName)

	if err != nil {
		return "", err
	}

	return string(buf), nil
}

//ModTime implement `Reader.ModTime`
func (r *FsReader) ModTime(ctx context.Context) (int64, error) {
	fi, err := os.Stat(r.fileName)
	if err != nil {
		return 0, err
	}
	return fi.ModTime().UnixNano(), nil
}
