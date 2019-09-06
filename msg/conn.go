package msg

import (
	"encoding/binary"
	"fmt"
	"penet/conn"
	"time"
)


//从conn中读取byte并反序列化为message
//将message序列化为bytes并通过conn发送

func WriteMsg(conn conn.PConn, m Message) (err error){
	var buffer []byte

	//pack消息
	if buffer, err = Pack(m); err != nil{
		return
	}

	msgLen := len(buffer)

	//先按照大端法写入信息长度 int
	if err = binary.Write(conn, binary.BigEndian, &msgLen); err != nil{
		return
	}

	conn.Info("Waiting to write %d bytes", msgLen)

	//写入信息数据并设置超时
	// 5s写超时的时间
    if err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil{
    	return
	}

    var n int

    if n , err = conn.Write(buffer); err != nil{
    	return
	}

    if n != msgLen{
    	err = fmt.Errorf("Expected %d bytes to write, but only write %d bytes ", msgLen, n)
    	return
	}
    return
}

func readMsg(conn conn.PConn, msgIn Message) (m Message, err error){
    //读取信息长度
    var msgLen int

    //大端法读取信息长度
    if err = binary.Read(conn, binary.BigEndian, &msgLen); err != nil{
    	return
	}

    conn.Info("Waiting to read %d bytes", msgLen)

    //todo 如何避免重复分配缓冲区
    //分配缓冲区
    buffer := make([]byte, msgLen)

    //设置5s读取超时

    if err = conn.SetReadDeadline(time.Now().Add(5*time.Second)); err != nil{
    	return
	}

    var n int

    if n, err = conn.Read(buffer); err != nil{
    	return
	}

    if n != msgLen{
    	err = fmt.Errorf("Expected %d bytes to read, but only read %d bytes ", msgLen, n)
    	return
	}

    if msgIn == nil{
    	return UnPack(buffer)
	}

    return UnPackInto(buffer, msgIn)
}


func ReadMsg(conn conn.PConn) (Message, error){
    return readMsg(conn, nil)
}

func ReadMsgInto(conn conn.PConn, msgIn Message)(Message, error){
	return readMsg(conn, msgIn)
}

