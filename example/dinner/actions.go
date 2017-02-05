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
	for {
		p.think()
		fmt.Printf("[%v] is hungry\n", p.ID)
		ld.In(ticket{})
		ld.In(chopstick(i))
		ld.In(chopstick((i + 1) % num))
		p.eat()
		ld.Out(chopstick(i))
		ld.Out(chopstick((i + 1) % num))
		ld.Out(ticket{})
	}
}

func (p philosopher) think() {
	fmt.Printf("[%v] is thinking\n", p.ID)
	time.Sleep(time.Duration(rand.Int31n(2000)) * time.Millisecond)
	fmt.Printf("[%v] has finished thinking\n", p.ID)
}
func (p philosopher) eat() {
	fmt.Printf("[%v] is eating\n", p.ID)
	time.Sleep(time.Duration(rand.Int31n(2000)) * time.Millisecond)
	fmt.Printf("[%v] has finished eating\n", p.ID)
}
