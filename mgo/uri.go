package mgo

import (
	"net/url"
	"strings"

	"github.com/yaitoo/sparrow/types"
)

//ConnectionString Mongo连接信息
//mongodb://[username:password@]host1[:port1][,...hostN[:portN]]][/[database][?options]]
type URI struct {
	//Hosts 主机列表
	Hosts []string
	//Database 数据库
	Database string
	//AuthDB 验证数据库
	AuthDB string
	//Login 登陆账号
	Login string
	//Passwd 登陆密码
	Passwd string
	//ReplicaSet 复制分片
	ReplicaSet string
}

//NewURI 创建新都连接池
func NewURI(host, database, authdb, replicaSet, login, passwd string) *URI {
	uri := &URI{}
	uri.Hosts = strings.Split(strings.TrimSpace(host), ",")
	uri.Database = database

	if types.IsEmpty(authdb) {
		uri.AuthDB = database
	} else {
		uri.AuthDB = authdb
	}
	uri.Login = login
	uri.Passwd = passwd

	uri.ReplicaSet = replicaSet

	return uri
}

func (u *URI) String() string {
	uri := "mongodb://"
	hasOption := false

	if types.IsNotEmpty(u.Login) && types.IsNotEmpty(u.Passwd) {
		uri += url.QueryEscape(u.Login) + ":" + url.QueryEscape(u.Passwd) + "@"
	}

	uri += strings.Join(u.Hosts, ",")

	if types.IsNotEmpty(u.Database) {
		uri += "/" + u.Database
	}

	if types.IsNotEmpty(u.Login) && types.IsNotEmpty(u.Passwd) && types.IsNotEmpty(u.AuthDB) {
		hasOption = true
		uri += "?authSource=" + u.AuthDB
	}

	if types.IsNotEmpty(u.ReplicaSet) {
		if hasOption {
			uri += "&replicaSet=" + u.ReplicaSet
		} else {
			uri += "?replicaSet=" + u.ReplicaSet
		}
	}

	return uri
}
