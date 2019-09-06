package conn

import (
	"fmt"
	"io"
	"net"
	"penet/log"
	"sync"
)

//处理所有的连接


type PConn interface {
	//带前缀日志
	log.Logger

	//
	net.Conn

	//proxy | control
	SetType(typ string)

	//唯一的能够清楚的表示这个conn的信息
	Id() string
}


//
type LoggedConn struct {
	//带prefixes日志
	log.Logger

	//持有一个conn
	net.Conn

	//随机id
	id string

	//typ
	typ string
}


func NewLoggedConn(c net.Conn, prefixes ...string) *LoggedConn{
	lc := &LoggedConn{
		Logger: log.NewPrefixedLogger(prefixes...),
		Conn: c,
	}
	return lc
}

func (lc *LoggedConn) Id() string{
	return fmt.Sprintf("%s:%s", lc.typ, lc.id)
}

func (lc *LoggedConn) SetType(typ string){
	oldId := lc.Id()
	lc.typ = typ
	newId := lc.Id()

	//日志
	lc.Info("Connection renamed from %s to %s", oldId, newId)
}


func WrapConn(conn net.Conn, typ string) PConn{
	switch c := conn.(type) {
	case *LoggedConn:
		return c
	case *net.TCPConn:
		return NewLoggedConn(conn, typ)
	}
	return nil
}

//将两个conn连接起来转发数据
//c1收到的数据转发给c2,c2收到的数据转发给c1
//fromBytes是c1转发给c2的数据字节
//toBytes是c2转发给c1的数据字节

func Join(c1 PConn, c2 PConn) (fromBytes int64, toBytes int64){
	var wait sync.WaitGroup

	//todo 通用函数处理此类错误
	defer c1.Close()
	defer c2.Close()

	pipe := func(from PConn, to PConn, fromBytes *int64) {
		//执行完则done否则会死锁
		defer wait.Done()
		var err error
		*fromBytes, err = io.Copy(to, from)

		if err != nil{
			from.Warn("Copied %d bytes before fail with %v", *fromBytes, err)
		}//if

	}//pipe

	//等待两个携程执行完毕
	wait.Add(2)

	go pipe(c1, c2, &fromBytes)
	go pipe(c2, c2, &toBytes)

	//logging
	c1.Info("Copied %d bytes from %s, %d bytes from %s", fromBytes, c1.Id(), toBytes, c2.Id())
	return
}