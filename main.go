package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

func main() {
	helpPtr:= flag.Bool("help", false, "prints help")
	tickerPtr := flag.Bool("ticker", false, "runs ticker loop")
	sleepPtr:=flag.Bool("sleep", false, "runs loop with sleep")
	noGoPtr:=flag.Bool("nogo", false, "runs function that does not use go routines")
	goPtr:=flag.Bool("go", false, "runs function that uses go routines")
	inGoPtr:=flag.Bool("ingo", false, "runs function in a go routine")

	flag.Parse()

	if *helpPtr{
		fmt.Print(`usage: ./tickertest [flags]

TickerTest executes a long-running function. This is to test what happens if you call a function using a Ticker that takes longer than a tick. For comparison, it can also call the long-running function in a loop using time.Sleep.

The flags are:

	-ticker
			runs the long-running function in a ticker loop (default is true)
	-sleep
			runs the function in a loop with a time.Sleep
	-nogo
			runs a function that does not use go routines (default is true)
	-go
			runs a function that uses go routines
	-ingo
			runs a function in a go routine
	-help
			print help
`)
		return
	}

	// Every 100 milliseconds call the doSomethingWithGoRoutines() method, until we've done it 5 times.
	// The go routines in doSomethingWithGoRoutines() will take longer than 100 milliseconds each, so we'll be starting
	// another call to doSomethingWithGoRoutines() when some go routines from the first call are still running.
	fmt.Printf("app start time %s\n", time.Now().String())

	if *tickerPtr {
		ticker := time.NewTicker(100 * time.Millisecond)
		i := -1
		previousTick := time.Now()

		for t := range ticker.C {
			i++

			fmt.Printf("tick %d start time %s\n", i, t.String())
			timeSinceLastTick := float64(time.Since(previousTick).Nanoseconds() / 1e6)
			fmt.Printf("tick %d time since last tick %fms\n", i, timeSinceLastTick)
			previousTick = t

			if i == 5 {
				break
			}

			// Because we execute code which blocks, there will be no new ticks until
			// the code returns. This means we won't start backing up multiple executions of the called function.
			// If we did `go doSomethingWithGoRoutines(i)`, then we would start overlapping executions.
			if *noGoPtr {
				doSomething(i)
			}

			if *goPtr {
				doSomethingWithGoRoutines(i)
			}

			if *inGoPtr {
				go doSomething(i)
			}

			fmt.Printf("tick %d finish time %s\n", i, time.Now().String())
		}
	}

	if *sleepPtr {
		j := -1
		for {
			fmt.Printf("loop %d start time %s\n", j+1, time.Now().String())

			j++
			if j == 5 {
				break
			}

			if *noGoPtr {
				doSomething(j)
			}

			if *goPtr {
				doSomethingWithGoRoutines(j)
			}

			if *inGoPtr {
				go doSomething(j)
			}

			fmt.Println("sleeping for 100 milliseconds")
			time.Sleep(100 * time.Millisecond)

			fmt.Printf("loop %d end time %s\n", j, time.Now().String())
		}
	}
	fmt.Println("Finished")
}

func doSomethingWithGoRoutines(i int) {
	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(10)
	total := 0

	// Create 10 go routines and print out when they start and finish.
	// Do a thing that takes longer than 100 milliseconds.
	for j := 0; j < 10; j++ {
		go func(i, j int) {
			var lTotal int
			defer wg.Done()
			//fmt.Printf("doSomethingWithGoRoutines/iteration %d/%d started %s\n", i, j, time.Now().String())
			for k := 0; k < 100000000; k++ {
				lTotal++
			}
			total += lTotal
			//fmt.Printf("doSomethingWithGoRoutines/iteration %d/%d finished %s with total %d\n", i, j, time.Now().String(), total)
		}(i, j)
	}

	wg.Wait()

	duration := float64(time.Since(startTime).Nanoseconds() / 1e6)
	fmt.Printf("doSomethingWithGoRoutines %d finished %s (duration %fms)\n", i,  time.Now().String(),duration)
}

func doSomething(i int) {
	startTime := time.Now()

	total := 0

	// Print out start and end times of each loop.
	// Do a thing that takes longer than 100 milliseconds.
	for j := 0; j < 10; j++ {

		var lTotal int

		//fmt.Printf("doSomething/iteration %d/%d started %s\n", i, j, time.Now().String())
		for k := 0; k < 100000000; k++ {
			lTotal++
		}
		total += lTotal
		//fmt.Printf("doSomething/iteration %d/%d finished %s with total %d\n", i, j, time.Now().String(), total)
	}

	duration := float64(time.Since(startTime).Nanoseconds() / 1e6)
	fmt.Printf("doSomething %d finished %s (duration %fms)\n", i,  time.Now().String(),duration)
}