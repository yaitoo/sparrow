package mgo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

//Session mongo進程
type Session struct {
	Client *mongo.Client
	ctx    context.Context
	status int
	retry  int
	closed bool
	sleep  bool

	onSessionClosed func(*Session)
}

//Close 關閉連接， 並返回連接池
func (s *Session) Close() {
	if s != nil && s.onSessionClosed != nil {
		go s.onSessionClosed(s)
	}
}

//connect 尝试连接mongodb
func (s *Session) connect() error {

	if err := s.Client.Connect(s.ctx); err != nil {
		s.status = 0
		return err
	}

	s.status = 1
	go s.keepalive()
	return nil

}

//disconnect
func (s *Session) disconnect() error {
	logger.Printf("%v disconnect", s)
	s.closed = true
	return s.Client.Disconnect(s.ctx)
}

//keepalive 保持连接
func (s *Session) keepalive() {
	go func() {
		t := time.NewTicker(1 * time.Second)
		for {
			<-t.C
			if s.closed {
				break
			}

			if s.sleep == false {
				continue
			}

			ctx, cancel := context.WithTimeout(s.ctx, time.Second)
			if err := s.Client.Ping(ctx, nil); err != nil {

				s.retry++
				s.status = 0
				s.closed = true

				logger.Warnf("Ping err %s, retry %v", err, s.retry)

				if s.retry > 1 {
					s.disconnect()
				}
			} else {
				s.status = 1
				s.retry = 0
				//logger.Printf("%v keepalive...\n", s)
			}
			cancel()
		}

	}()
}

func wrapAsSession(ctx context.Context, c *mongo.Client) *Session {
	s := &Session{}
	s.ctx = ctx
	s.Client = c
	s.status = 0

	return s
}
