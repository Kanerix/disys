package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Fork struct {
	id                         int
	used                       bool
	request, response, release chan bool
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
			id:       i,
			request:  make(chan bool, 1),
			response: make(chan bool, 1),
			release:  make(chan bool, 1),
		}
		go forks[i].Serve()
	}

	philosophers := make([]Philosopher, philosopherCount)
	for i := range philosophers {
		philosophers[i] = Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%philosopherCount],
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

// A monitor goroutine that stands for communication with the fork.
//
// If a request is successful we set the `used` flag to `true`
// and communicate back to the requester that the fork was taken.
// Otherwise, we communicate back that the fork was unavailable.
//
// If a fork is released we set the `used` flag to `false`.
func (f Fork) Serve() {
	for {
		select {
		case <-f.request:
			// We communicate back to the requester if the request
			// was successful returning the `true` if the fork was taken
			// and `false` if the fork was unavailable.
			if !f.used {
				// This is not a race condition since we are only
				// reading and writing to `f.used` on one thread.
				f.used = true

				fmt.Println("Fork", f.id, "is now being used")
				f.response <- true
			} else {
				fmt.Println("Fork", f.id, "is unavailable")
				f.response <- false
			}
		case <-f.release:
			fmt.Println("Fork", f.id, "is now released")
			f.used = false
		}
	}
}

// Tries to grap the fork by sending any message to the request channel.
//
// This will return the response from the fork's monitor goroutine.
// If the fork is available we return `true` otherwise `false`.
// If communication fails we return `false` as well.
func (f Fork) TryGrap() bool {
	select {
	case f.request <- true:
		return <-f.response
	default:
		return false
	}
}

// Releases the fork making it available for other philosophers.
//
// fork's release channel. This will make the fork
// available for other philosophers to use.
// Release a fork so other philosophers can use it.
func (f Fork) Release() {
	f.release <- true
}

// Start the philosopher's routine.
//
// The philosopher will try to dine until he has eaten 3 meals.
func (p *Philosopher) Philosophize(wg *sync.WaitGroup) {
	defer wg.Done()
	for p.meals < 3 {
		p.Dine()
	}
}

// The philosopher's dinning routine.
//
// The philosopher will try to grap the right fork first.
// If he can't grap the right fork he will try again.
// If he can't grap the left fork he will release the right fork
// and try again.
// If he can grap both forks he will eat and release the forks.
func (p *Philosopher) Dine() {
	for {
		p.Think()

		hasRightFork := p.rightFork.TryGrap()
		if !hasRightFork {
			continue
		}

		hasLeftFork := p.leftFork.TryGrap()
		if !hasLeftFork {
			// Here we make sure the program doesn't get stuck in a deadlock
			// since we drop the right fork if we can't get the left one
			// making it available for other philosophers to pick up.
			p.rightFork.Release()
			continue
		}

		p.Eat()

		// Make sure to release the forks after eating so other
		// philosophers can use them again also avoiding deadlocks.
		p.rightFork.Release()
		p.leftFork.Release()
		break
	}
}

// The philosopher's thinking routine.
//
// The philosopher will think for a random amount of time.
func (p *Philosopher) Think() {
	fmt.Println("Philosopher", p.id, "is thinking")
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)+200))
}

// The philosopher's eating routine.
//
// The philosopher will eat for a random amount of time.
func (p *Philosopher) Eat() {
	fmt.Println("Philosopher", p.id, "is eating his", p.meals+1, "meal")
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(800)+200))
	p.meals++
}
