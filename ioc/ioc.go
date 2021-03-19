// Copyright 2019 Sparrow. All rights reserved.

// Package ioc implements inversion of control containers and the dependency injection pattern.
// see https://martinfowler.com/articles/injection.html
package ioc

import (
	"context"
	"errors"
	"sync"

	"github.com/yaitoo/sparrow/log"
)

//ErrNoRegistered 物件/實例未被註冊
var ErrNoRegistered = errors.New("ioc: type/instance is not registered")

//ErrHasRegistered 物件/實例已經被註冊
var ErrHasRegistered = errors.New("ioc: type/instance has already been registered")

var (
	createLocker   sync.RWMutex
	instanceLocker sync.RWMutex
	creates        = make(map[string]CreateNew)
	instances      = make(map[string]interface{})
	logger         = log.NewLogger("ioc")
)

//Register 註冊物件
func Register(ctx context.Context, name string, createNew CreateNew) error {
	createLocker.Lock()
	defer createLocker.Unlock()
	_, ok := creates[name]
	if ok {
		logger.Warn("type [" + name + "] has been registered.")
		return ErrHasRegistered
	}

	creates[name] = createNew

	return nil

}

//RegisterInstance 註冊全局的物件實例
func RegisterInstance(ctx context.Context, name string, instance interface{}) error {
	instanceLocker.Lock()
	defer instanceLocker.Unlock()
	_, ok := instances[name]
	if ok {
		logger.Warn("instance [" + name + "] has been registered.")
		return ErrHasRegistered
	}

	instances[name] = instance

	return nil

}

//Resolve  創建指定的名稱的物件新實例
func Resolve(ctx context.Context, name string) (interface{}, error) {
	createLocker.RLock()
	defer createLocker.RUnlock()

	fn, ok := creates[name]
	if ok {
		return fn(ctx), nil
	}

	logger.Warn("type name [" + name + "] is not registered.")
	return nil, ErrNoRegistered

}

//ResolveInstance 獲取全局的物件實例
func ResolveInstance(ctx context.Context, name string) (interface{}, error) {
	instanceLocker.RLock()
	defer instanceLocker.RUnlock()

	inst, ok := instances[name]
	if ok {
		return inst, nil
	}

	logger.Warn("instance name [" + name + "] is not registered.")
	return nil, ErrNoRegistered
}

//CreateNew 購造新物件實例
type CreateNew func(ctx context.Context) interface{}
