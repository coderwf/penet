package server

import (
	"net"
	"penet/conn"
	"penet/msg"
	"time"
)

func handleClientConn(c conn.PConn){
	//处理客户端连接请求
	var err error
	var m msg.Message
	err = c.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil{
		//todo panic ?
		panic(err)
	}

    m, err = msg.ReadMsg(c)
    if err != nil{
        //todo do what
        return
	}

	switch rawMsg := m.(type) {
	case *msg.Auth:
		NewControl(rawMsg, c)
	case *msg.ProxyReq:
		NewProxy(rawMsg, c)
	}
}


func handlePubConn(c net.Conn){
    c.LocalAddr().String()
}