package conn

import (
	"net"
	"penet/log"
)


type PListener struct {

	//监听到的所有连接
	Conns chan PConn

	//监听的地址
	net.Addr
}


//开始监听某个端口并返回PListener
func Listen(typ string, address string) (pListener *PListener, err error){
    l , err := net.Listen("tcp", address)
    if err != nil{
    	return
	}

    pListener = &PListener{
    	Conns:make(chan PConn, 5),
        Addr: l.Addr(),
	}

    var c net.Conn
    var e error
    accept := func() {
		for{
			c, e = l.Accept()
			//todo close
			if e != nil{
				//log
				log.Warn("Listen %s accept fail with %v", l.Addr().String(), e)
				continue
			}

			//将c放入管道
			log.Info("New connection from %s", c.RemoteAddr().String())
			pListener.Conns <- WrapConn(c, typ)
		}//for
	}

    go accept()
    return
}
