package server

import (
	"penet/conn"
	"penet/log"
)

//监听client到来的连接,可能为control或者proxy
func tunnelListener(network string, address string){
	listener, err := conn.Listen("tunnel", address)
	if err != nil{
		panic(err)
	}

	//步处理连接请求
	//todo 怎么关闭Conns通道
	for c := range listener.Conns{
		go handleClientConn(c)
	}//for

}


func Main(){
	log.LogToStdout("DEBUG")
    tunnelListener("tcp", ":6666")
}
