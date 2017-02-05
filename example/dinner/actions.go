package main

import (
	"fmt"
	"math/rand"
	"time"
)

type philosopher struct {
	ID int
}

func phil(i int) {
	p := philosopher{i}
	fmt.Printf("Philosopher %v is born\n", p.ID)
}

func (p philosopher) think() {
	fmt.Printf("[%v] is thinking", p.ID)
	time.Sleep(time.Duration(rand.Int31n(5000)) * time.Millisecond)
}
func (p philosopher) eat() {
	fmt.Printf("[%v] is eating", p.ID)
	time.Sleep(time.Duration(rand.Int31n(5000)) * time.Millisecond)
}
