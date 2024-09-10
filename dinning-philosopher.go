package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Fork struct {
	id      int
	request chan bool
	release chan bool
}

type Philosopher struct {
	id                  int
	rightFork, leftFork Fork
	meals               int
}

const philosopherCount = 5

func main() {
	forks := make([]Fork, philosopherCount)
	for i := range forks {
		forks[i] = Fork{
			id:      i,
			request: make(chan bool, 1),
			release: make(chan bool, 1),
		}
		go forks[i].Serve()
	}

	philosophers := make([]*Philosopher, philosopherCount)
	for i := 0; i < 5; i++ {
		philosophers[i] = &Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%5],
			meals:     0,
		}
	}

	// The wait group is a way to wait for all the philosophers to finish their
	// meals. All 5 philosophers will be dinning at the same time and will first
	// leave when they have eaten exactly 3 times.
	wg := sync.WaitGroup{}

	for _, philphilosopher := range philosophers {
		wg.Add(1)
		go philphilosopher.Philosophize(&wg)
	}

	// We use the wait group to wait for all the philosophers to finish their meals.
	// This makes sure that the main function does not exit before all the philosophers
	// are done eating.
	wg.Wait()
}

func (f Fork) Serve() {
	for {
		// We wait for a philosopher to request the fork.
		<-f.request
		// We then let the philosopher know that the fork is ready to be used.
		f.request <- true
		// We then wait for the philosopher to release the fork again.
		<-f.release
	}
}

func (p *Philosopher) Philosophize(wg *sync.WaitGroup) {
	defer wg.Done()
	for p.meals < 3 {
		p.Dine()
	}
	fmt.Println("Philosopher", p.id, "is done eating", p.meals, "meals and left the table")
}

func (p *Philosopher) Dine() {
	p.Think()

	// Here we request the forks and tell it we want to use it.
	// This will be picked up by the monitor group (`fork.Serve()`)
	p.rightFork.request <- true
	// We wait for the monitor group to tell us that we can use the fork.
	<-p.rightFork.request

	// The same logic goes for the left fork aswell.
	p.leftFork.request <- true
	<-p.leftFork.request

	// We know have both forks and can eat.
	p.Eat()

	// After we are done eating we release the forks.
	p.rightFork.release <- true
	p.leftFork.release <- true
}

func (p *Philosopher) Think() {
	fmt.Println("Philosopher", p.id, "is thinking")
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+200))
}

func (p *Philosopher) Eat() {
	fmt.Println("Philosopher", p.id, "is eating meal", p.meals+1)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+200))
	p.meals++
}
