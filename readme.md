# Go TickerTest

Experiment to see how Go's time.Ticker works with go routines.

It turns out if you call code normally on each tick, the next tick won't happen even if it's scheduled.

It doesn't matter if you have go routines inside that code.

When the code finishes, the next tick can then occur.

**Example**
```
// New ticks won't happen, even if they are scheduled, until the called code finishes.
// This avoids multiple concurrent executions of doSomething().
for t := range ticker.C {
	doSomething()
}
```

If you look at the timings, the next tick starts slightly before the code executed during the last tick ends.

If you were to run the code in a go routine, then each subsequent tick would trigger a new call to the code in a go routine.

**Example**
```
// New ticks will happen, creating multiple go routines.
for t := range ticker.C {
	go doSomething(i)
}
```

## How To Run

`go run main.go -help` to see what flags there are.