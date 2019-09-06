package server

import (
	"fmt"
	"penet/log"
	"sync"
)

//
var controlRegistry *ControlRegistry


func init(){
    controlRegistry = &ControlRegistry{
    	table: make(map[string] *Control),
	}
}

type ControlRegistry struct {
	mu sync.Mutex
	table map[string] *Control
}


func (cr *ControlRegistry) Register(ClientId string, ctl *Control) {
    //todo replace
	cr.mu.Lock()
	defer cr.mu.Unlock()

	var ok bool
	if _, ok = cr.table[ClientId]; ok{
		log.Warn("Control %s already in registry ", ClientId)
		return
	}
	cr.table[ClientId] = ctl
	return
}

func (cr *ControlRegistry) Get(ClientId string) (ctl *Control, err error){
	cr.mu.Lock()
	defer cr.mu.Unlock()

	var ok bool
	if ctl, ok = cr.table[ClientId]; !ok{
		err = fmt.Errorf("No control %s found in registry ", ClientId)
	}
	return
}

func (cr *ControlRegistry) Del(ClientId string) (err error){
    //todo 简单删除后资源怎么处理
    cr.mu.Lock()
    defer cr.mu.Unlock()

    if _, ok := cr.table[ClientId]; !ok{
    	err = fmt.Errorf("No control %s found in registry ", ClientId)
    	return
	}

    delete(cr.table, ClientId)
    return
}


func (cr *ControlRegistry) Len() int{
	cr.mu.Lock()
	defer cr.mu.Unlock()

	return len(cr.table)
}

