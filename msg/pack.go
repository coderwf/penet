package msg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

//打包,解包消息
//序列化,反序列化

type EnvLope struct {
	//消息类型,动态new该类型对应的消息结构体
	Typ string

	//消息序列化后的bytes
	Payload json.RawMessage
}


func Pack(m Message) (buffer []byte, err error){
	return json.Marshal(struct {
		typ string
		payload interface{}
	}{
		typ: reflect.TypeOf(m).Elem().Name(),
		payload:m,
	})
}

func unpack(buffer []byte, msgIn Message) (m Message, err error){
    var env EnvLope
    var ok bool
    if err = json.Unmarshal(buffer, &env); err != nil{
    	return
	}

    if msgIn == nil{
		m, ok = NewMsg(env.Typ)

		if !ok{
			err = fmt.Errorf("Message type %s not found ", env.Typ)
			return
		}//if
	}else{
		m = msgIn
	}//else

    err = json.Unmarshal(env.Payload, m)
    return
}


//将消息数据放入指定的msgIn中
func UnPackInto(buffer []byte, msgIn Message) (Message, error){
    return unpack(buffer, msgIn)
}


func UnPack(buffer []byte) (Message, error){
	return unpack(buffer, nil)
}