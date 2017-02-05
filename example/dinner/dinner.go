package main

import (
	"github.com/owulveryck/go-linda"
	"time"
)

var ld *linda.Linda

func main() {
	var done = make(chan struct{}, 0)
	ld = tupleSpace()
	for i := 0; i < 5; i++ {
		ld.Eval([]interface{}{phil, i})
	}
	go func() {
		time.Sleep(10 * time.Second)
		done <- struct{}{}
	}()
	<-done
}

func tupleSpace() *linda.Linda {
	input := make(<-chan interface{})
	output := make(chan<- interface{})
	ld := &linda.Linda{
		Input:  input,
		Output: output,
	}
	return ld
}
