package util

import "sync"

type ShutDown struct {
	mu sync.Mutex
	begin chan bool
	hasBegin bool
}


func NewShutDown() *ShutDown{
    return &ShutDown{
    	begin: make(chan bool),
    	hasBegin: false,
	}//return
}


func (s *ShutDown) WaitBegin(){
	<- s.begin
}


func (s *ShutDown) Begin(){
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.hasBegin{
		return
	}else{
		close(s.begin)
		s.hasBegin = true
	}
}
