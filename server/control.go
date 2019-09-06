package server

import (
	"fmt"
	"ngrok/util"
	"penet/conn"
	"penet/log"
	"penet/msg"
	"time"
)

//auth的control连接
type Control struct {
    //
	//logger
	log.Logger

	//控制连接
	conn conn.PConn

	//每一条控制连接的唯一标志
	//proxy携带这个ClientId表示为此control下的代理
	ClientId string

	//所有的来自client的代理连接
	proxies chan conn.PConn

	//最后一次接受到数据时间
	lastAccessTime time.Time


	//一个携程负责读取消息放入in队列
	//一个携程负责将out队列中的消息通过conn发送
	//一个携程负责处理in队列中的消息

	//ctl关闭的时候要关闭chan否则会死锁

	in chan msg.Message
	out chan msg.Message

    shutdown *util.Shutdown
}



//外部发起的请求
func pubListener() string{

	//随机监听一个端口
	listener, err := conn.Listen("public", "0.0.0.0:0")
	if err != nil{
		return ""
	}

	log.Info("Listener pub at: %s", listener.Addr.String())

	//开启携程异步处理连接请求
	go func() {
		//todo 怎么关闭Conns通道

		for c := range listener.Conns{
			go handlePubConn(c)
		}//for
	}()

	return listener.Addr.String()
}


func NewControl(auth *msg.Auth, c conn.PConn){
    //todo 权限验证
    var err error
    //todo 错误处理
    if auth.Token != "123456"{
    	err = msg.WriteMsg(c, &msg.AuthResp{
    		Error:"Auth Failed, Please check token",
		})
    	log.Info("%v", err)
	}//

	//权限验证通过则开始监听public请求
    listenAddr := pubListener()
    if listenAddr == ""{
		return
	}//if

	ClientId := fmt.Sprintf("%s-%s", util.RandId(5), listenAddr)
    err = msg.WriteMsg(c, &msg.AuthResp{
    	ClientId:ClientId,
    	RemoteAddr:listenAddr,
	})

    c.SetType("ctl")

    ctl := &Control{
        Logger: log.NewPrefixedLogger("Ctl", ClientId),
        conn: c,
        ClientId: ClientId,
        proxies: make(chan conn.PConn, 20),
        lastAccessTime: time.Now(),

        in:make(chan msg.Message, 5),
        out: make(chan msg.Message, 5),

        shutdown: util.NewShutdown(),

	}
    controlRegistry.Register(ClientId, ctl)
}

func NewProxy(req *msg.ProxyReq, proxyConn conn.PConn) {
    ctl , err := controlRegistry.Get(req.ClientId)
    if err != nil{
    	log.Info("No control found for ClientId: %s", req.ClientId)
		return
	}

    ctl.RegisterProxy(proxyConn)
}


func (ctl *Control) RegisterProxy(proxy conn.PConn){
	select {
	case ctl.proxies <- proxy:
        	ctl.Info("Register proxy %s", proxy.Id())
	default:
		//proxies is full ,discard
		ctl.Info("Proxy is full, discard proxy %s", proxy.Id())
	}
}

func (ctl *Control) GetProxy() (proxy conn.PConn, err error){
	var ok bool

	select{
	//拿到proxy则返回
	case proxy, ok = <- ctl.proxies:
		if !ok{
			//proxies通道关闭,直接返回错误
			err = fmt.Errorf("No avaiable proxy , closing ")
			return
		}
		return
	default:
		//没有拿到可用的proxy则需要通知客户端发起proxy请求
		//todo 如何保持proxy的复用性,保持有一定量的proxy待命
		ctl.out <- &msg.NewProxy{}
		ctl.out <- &msg.NewProxy{}
	}

	//再次重新获取proxy,如果超时则表示客户端无法发起proxy连接,此时结束ctl

	select{
	case proxy, ok = <- ctl.proxies:
		if !ok{
			//proxies通道关闭,直接返回错误
			err = fmt.Errorf("No avaiable proxy , closing ")
			return
		}
		return
	case <- time.After(10 * time.Second):
		err = fmt.Errorf("Client cant start proxy, closing ")
		return
	}

}


func (ctl *Control) reader(){

	defer ctl.shutdown.Begin()

	var m msg.Message
	var err error

    for{
		m, err = msg.ReadMsg(ctl.conn)
		if err != nil{
		    //todo close
			return
		}//if

		ctl.in <- m
	}//for
}

func (ctl *Control) writer(){

	defer ctl.shutdown.Begin()

	var err error

    for m := range ctl.out{
    	err = msg.WriteMsg(ctl.conn, m)

    	if err != nil{
			//todo close
			return
		}
	}//for
}


func (ctl *Control) manager(){
    //读取消息并处理,判断超时

    defer ctl.shutdown.Begin()

    var m msg.Message
    var err error

    ticker := time.NewTicker(10 * time.Second)

    //读取消息并判断超时
    for{
		select {
    	case m = <- ctl.in:
    		err = ctl.process(m)
    		if err != nil{
    			//todo log and close
    			return
			}
    	case <- ticker.C:
    		//判断超时
    		if time.Now().Sub(ctl.lastAccessTime) > time.Second * 20{
    			//heartbeat for 20 s
    			ctl.Info("Time out for 20s, closing")

    			//todo close
				return
    			//todo close

			}//if
		}//select
	}//for
}


func (ctl *Control) process(m msg.Message) (err error){
	switch m.(type) {
	case *msg.Ping:
		//更新最近通信时间
		ctl.lastAccessTime = time.Now()
		//发送pong心跳包
		ctl.out <- &msg.Pong{}

	}//switch
	return
}


//等待control关闭,执行销毁chan,关闭conn等操作
func (ctl *Control) stopper(){
    //等待结束
    ctl.shutdown.WaitBegin()

    //结束则执行清理工作

    //chan
    close(ctl.in)
    close(ctl.out)

    //conn
    ctl.conn.Close()
}
