package msg

import (
    "reflect"
)

var typeMap map[string] reflect.Type


func init(){
    //初始化typeMap
    typeMap = make(map[string] reflect.Type)

    t := func(obj interface{}) {
        typ := reflect.TypeOf(obj).Elem()
        typeMap[typ.Name()] = typ
    }

    //注册type
    t((*Auth)(nil))
    t((*AuthResp)(nil))
    t((*NewProxy)(nil))
    t((*ProxyReq)(nil))
    t((*StartProxy)(nil))
    t((*Ping)(nil))
    t((*Pong)(nil))

}


func NewMsg(typ string) (msg Message, ok bool){
    msgT, ok := typeMap[typ]
    if !ok{
        return
    }

    msg = reflect.New(msgT).Elem().Interface().(Message)
    return
}

//所有传递的消息

type Message interface {}


//client发送给server请求一个control连接

type Auth struct {
	//账号密码验证
    User string
    Password string

    //或者通过token验证也可以
    Token string

}


//server回复client的Auth消息,并通知client访问server哪个地址
//可以代理到client的本地服务

type AuthResp struct {
	//client发送代理连接请求时携带ClientId则可以将此proxy将和control对应
    ClientId string

    //告诉client用户可通过此地址访问到此control的代理进而访问到本地服务
    RemoteAddr string

    //可能认证失败,否则为空字符串
    Error string
}


//server通过ctl告诉client该发起一个代理请求连接
type NewProxy struct {

}

//client通过连接listener并发送ProxyReq消息表示此连接为proxy连接
// (可能为control连接请求或者proxy连接请求)

type ProxyReq struct {
	//发送ClientId告诉server这个proxy归某个ctl
	ClientId string
}

//server通过proxy连接告诉proxy连接开始代理
type StartProxy struct {

}


//heartbeat

type Ping struct {

}

type Pong struct {

}

