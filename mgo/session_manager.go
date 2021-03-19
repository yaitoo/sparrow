package mgo

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/yaitoo/sparrow/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ErrWaitTimeout 获取连接池对象超时
var ErrWaitTimeout = errors.New("mgo: wait timeout")
var logger = log.NewLogger("mgo")

//SessionManager session管理器
type SessionManager interface {
	//Get 获取一个Mongodb链接
	Get() (*Session, error)
}

type sessionManager struct {
	ctx          context.Context
	conns        chan *Session
	maxIdleConns int
	maxOpenConns int32
	waitTimeout  time.Duration
	opts         []*options.ClientOptions
	total        int32
}

//NewSessionManager 初始化Session管理器实例
//maxIdleConns 最大空闲连接数
//waitTimeout  等待可用连接超时时间
func NewSessionManager(ctx context.Context, maxIdleConns, maxOpenConns int, waitTimeout time.Duration, opts ...*options.ClientOptions) (SessionManager, error) {
	sm := &sessionManager{}
	sm.ctx = ctx
	sm.maxIdleConns = maxIdleConns
	sm.maxOpenConns = int32(maxOpenConns)
	sm.conns = make(chan *Session, maxIdleConns)
	sm.opts = opts
	sm.waitTimeout = waitTimeout

	for i := 0; i < maxIdleConns; i++ {
		conn, err := sm.createNew()
		if err == nil {
			break
		}

		go atomic.AddInt32(&sm.total, 1)
		sm.pushBack(conn)

	}

	return sm, nil
}

func (sm *sessionManager) Get() (*Session, error) {
next:
	select {
	case conn := <-sm.conns:
		conn.sleep = false
		if conn.closed {
			goto next
		}

		if conn.status == 0 {
			//go sm.pushBack(conn)
			goto next
		}
		sm.printStatus("Get")

		return conn, nil

	case <-time.After(sm.waitTimeout):
		if atomic.LoadInt32(&sm.total) > sm.maxOpenConns {
			goto next
		}
		conn, err := sm.createNew()
		if err != nil {
			return nil, err
		}
		conn.sleep = false

		go atomic.AddInt32(&sm.total, 1)

		sm.printStatus("Get timeout")

		return conn, nil
	}
}

func (sm *sessionManager) createNew() (*Session, error) {
	c, err := mongo.NewClient(sm.opts...)
	if err != nil {
		return nil, err
	}

	conn := wrapAsSession(sm.ctx, c)

	if err := conn.connect(); err != nil {
		go conn.disconnect()
		return nil, err
	}

	conn.onSessionClosed = sm.pushBack

	return conn, nil

}

func (sm *sessionManager) pushBack(s *Session) {
	if s != nil {
		select {
		case sm.conns <- s:
			s.sleep = true
			sm.printStatus("Pushback")
		case <-time.After(sm.waitTimeout):
			go atomic.AddInt32(&sm.total, -1)
			sm.printStatus("Pushback timeout")
			s.sleep = false
			s.disconnect()
		}
	}
}

func (sm *sessionManager) printStatus(topic string) {
	logger.Printf("%s: total:%v size:%v brust:%v sleep: %v \n", topic, sm.total, sm.maxIdleConns, sm.maxOpenConns, len(sm.conns))
}
