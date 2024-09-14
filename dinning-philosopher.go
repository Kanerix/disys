package main

import (
	"fmt"
	"sync"
)

type Fork struct {
	id      int
	used    bool
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

	wg := sync.WaitGroup{}

	for _, philphilosopher := range philosophers {
		wg.Add(1)
		go philphilosopher.Philosophize(&wg)
	}

	wg.Wait()
}

func (f Fork) Serve() {
	for {
		select {
		case <-f.request:
			if !f.used {
				fmt.Println("Fork", f.id, "is now being used")
				f.used = true
			}
			f.request <- f.used
		case <-f.release:
			fmt.Println("Fork", f.id, "is now released")
			f.used = false
		}
	}
}

func (f Fork) Grap() bool {
	select {
	case f.request <- true:
		return <-f.request
	default:
		return false
	}
}

func (f Fork) Release() {
	f.release <- true
}

func (p *Philosopher) Philosophize(wg *sync.WaitGroup) {
	defer wg.Done()
	for p.meals < 3 {
		p.Dine()
	}
}

func (p *Philosopher) Dine() {
	for {
		p.Think()

		hasRightFork := p.rightFork.Grap()
		if !hasRightFork {
			continue
		}

		hasLeftFork := p.leftFork.Grap()
		if !hasLeftFork {
			p.rightFork.Release()
			continue
		}

		p.Eat()

		p.rightFork.Release()
		p.leftFork.Release()
		break
	}
}

func (p *Philosopher) Think() {
	fmt.Println("Philosopher", p.id, "is thinking")
	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+200))
}

func (p *Philosopher) Eat() {
	fmt.Println("Philosopher", p.id, "is eating")
	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+200))
	p.meals++
}
