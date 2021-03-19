package mux

import "sync"

//Router 路由表
type Router struct {
	sync.RWMutex
}
