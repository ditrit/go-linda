package main

import (
	"github.com/owulveryck/go-linda"
)

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
