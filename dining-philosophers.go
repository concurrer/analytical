/* Author: concurrer */

package main

/*
philosophers and forks are seated next to each other in such a way
that P0 has L0 to his left and L1 to his right
so P5 has L5 to his left and L0 to his right
That translates to Pi having Li to his left and L(i+1)%5 to his right
*/

import (
	"fmt"
	"sync"
	"syscall"
	"time"
)

const NumberOfPhilosophers = 5

type Philosopher struct {
	id, eatingCount           int
	gotLeftFork, gotRightFork bool
	eatChan                   chan bool
}
type DiningPhilosophers struct {
	Philosophers [NumberOfPhilosophers]Philosopher
	forks        [NumberOfPhilosophers]sync.Mutex
	wg           sync.WaitGroup
}

func (dp *DiningPhilosophers) pickLeftFork(philosopher int) {
	dp.forks[philosopher].Lock()
	dp.Philosophers[philosopher].gotLeftFork = true
	fmt.Printf("philosopher %d picked up leftFork\n", philosopher)
}
func (dp *DiningPhilosophers) pickRightFork(philosopher int) {
	dp.forks[(philosopher+1)%NumberOfPhilosophers].Lock()
	dp.Philosophers[philosopher].gotRightFork = true
	fmt.Printf("philosopher %d picked up rightFork\n", philosopher)
}

/*
gotBothForks() has:
 1. a timer to release blocked forks if both are not obtained in <100ms
 2. if both forks are obtained, sends a 'true' to the 'eatChan' for closing the eat loop
*/
func (dp *DiningPhilosophers) gotBothForks(philosopher int) {
	timer := time.NewTimer(100 * time.Millisecond)
	for {
		select {
		case <-timer.C:
			fmt.Printf("philosopher %d timed out waiting for both locks\n", philosopher)
			return
		default:
			if dp.Philosophers[philosopher].gotLeftFork && dp.Philosophers[philosopher].gotRightFork {
				fmt.Printf("philosopher %d got both locks\n", philosopher)
				dp.Philosophers[philosopher].eatChan <- true
				return
			} else {
				fmt.Printf("philosopher %d waiting for 2nd lock\n", philosopher)
			}
		}
	}
}

func (dp *DiningPhilosophers) putLeftFork(philosopher int) {
	dp.forks[philosopher].Unlock()
	dp.Philosophers[philosopher].gotLeftFork = false
	fmt.Printf("philosopher %d put down leftFork\n", philosopher)
}
func (dp *DiningPhilosophers) putRightFork(philosopher int) {
	dp.forks[(philosopher+1)%NumberOfPhilosophers].Unlock()
	dp.Philosophers[philosopher].gotRightFork = false
	fmt.Printf("philosopher %d put down rightFork\n", philosopher)
}

/*
eat():
 1. has a timer for 'eatChan' which in turn waits for 'gotBothForks()' so as to terminate select/timer
 2. completes the eat loop when a 'true' comes thru the 'eatChan' in select/case
*/
func (dp *DiningPhilosophers) eat(philosopher int) {
	for {
		fmt.Printf("philosopher %d is waiting to eat.\n", philosopher)

		// get the left fork -- this should be an independent call
		dp.wg.Add(1)
		go dp.pickLeftFork(philosopher)
		dp.wg.Done()

		// get the right fork -- this should be an independent call
		dp.wg.Add(1)
		go dp.pickRightFork(philosopher)
		dp.wg.Done()

		// check if both forks are obtained, if yes then send 'true' to eatChan for the next select to pick it up
		dp.gotBothForks(philosopher)

		select {
		case <-dp.Philosophers[philosopher].eatChan:
			// wait for both forks and then start eating
			fmt.Printf("philosopher %d got both forks and started eating\n", philosopher)
			dp.Philosophers[philosopher].eatingCount++ // exit if this is 5, that is when a philosopher eats 5 times
			dp.putRightFork(philosopher)
			dp.putLeftFork(philosopher)

			if dp.Philosophers[philosopher].eatingCount == 3 { // terminating condition for limited/unlimited eating cycles
				fmt.Printf("philosopher %d done eating with eatingCount: %d\n", philosopher, dp.Philosophers[philosopher].eatingCount)
				return
			} else {
				fmt.Printf("philosopher %d eatingCount: %d\n", philosopher, dp.Philosophers[philosopher].eatingCount)
			}
		case <-time.After(250 * time.Millisecond):
			// release the single lock to unblock others
			fmt.Printf("philosopher %d released single lock as timer expired\n", philosopher)
			if dp.Philosophers[philosopher].gotRightFork {
				dp.putRightFork(philosopher)
			}
			if dp.Philosophers[philosopher].gotLeftFork {
				dp.putLeftFork(philosopher)
			}
			break
		}

	}
}

func (dp *DiningPhilosophers) initPhilosophers() {
	for i := 0; i < NumberOfPhilosophers; i++ {
		dp.wg.Add(1)
		go func() {
			dp.eat(i)
			dp.wg.Done()
		}()
	}
}

func main() {
	/*
		Need a minimum 2 philosophers to form a circle
	*/
	if NumberOfPhilosophers < 2 { // condition for testcase, this is hardcoded as 5 for now
		fmt.Println("Number of philosophers MUST be at least 2")
		syscall.Exit(1)
	}
	dp := DiningPhilosophers{Philosophers: [5]Philosopher{
		Philosopher{eatChan: make(chan bool, 2)},
		Philosopher{eatChan: make(chan bool, 2)},
		Philosopher{eatChan: make(chan bool, 2)},
		Philosopher{eatChan: make(chan bool, 2)},
		Philosopher{eatChan: make(chan bool, 2)},
	}, forks: [5]sync.Mutex{}, wg: sync.WaitGroup{}}

	dp.initPhilosophers()
	dp.wg.Wait()
	fmt.Println("All philosophers done eating.")
}

/*
Debug output:
...
philosopher 4 done eating with eatingCount: 3
philosopher 3 picked up rightFork
philosopher 3 waiting for 2nd lock
philosopher 3 got both locks
philosopher 3 got both forks and started eating
philosopher 3 put down rightFork
philosopher 3 put down leftFork
philosopher 3 done eating with eatingCount: 3
All philosophers done eating.

Process finished with the exit code 0
*/
