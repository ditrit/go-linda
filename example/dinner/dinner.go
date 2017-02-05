package main

import (
	"github.com/owulveryck/go-linda"
	"time"
)

var ld *linda.Linda
var num = 5

type chopstick int
type ticket struct{}

func main() {
	var done = make(chan struct{}, 0)
	ld = tupleSpace()
	for i := 0; i < num; i++ {
		ld.Out(chopstick(i))
		ld.Eval([]interface{}{phil, i})
		if i < (num - 1) {
			ld.Out(ticket{})
		}
	}
	go func() {
		time.Sleep(30 * time.Second)
		done <- struct{}{}
	}()
	<-done
}

func tupleSpace() *linda.Linda {
	input := make(chan interface{}, 10)
	output := make(chan interface{}, 10)
	ld := &linda.Linda{
		Input:  input,
		Output: output,
	}
	go func() {
		for i := range output {
			input <- i
		}
	}()
	return ld
}
