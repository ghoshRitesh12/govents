## Introduction

Upon visiting this article, you might wonder how can Go have NodeJS like event contructs. Well, the reality is that Go has support for event driven architecture or Observer Pattern just like NodeJS, but not in a typical Node fashion. I encountered this when I was working on my last Go project named [Montre](https://riteshghosh.hashnode.dev/nodemon-clone-for-go), where I had to use the [fsnotify](https://github.com/fsnotify/fsnotify) package to watch for files or directory changes, which is obviously an asynchronous and event-driven task. I was fascinated when I implemented this feature in Go because it had the same vibe as a Node's `process.on(event, () => {})`.

Node exposes its event emitter API pretty neatly, as it's so intuitive to implement event-driven tasks, but that's not nearly the same with Go. In Go, we have to code something like the code snippet below to achieve somewhat Node like event drivenness.

```go
package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	sampleChannel := make(chan string)
	wg.Add(1)

    // this go routine acts as an asynchronous event listener
	go func() {
		for {
			select {
			case val, ok := <-sampleChannel:
				if !ok {
					log.Fatalln("channel is closed")
				}
				fmt.Println("received", val)
				wg.Done()
			}
		}
	}()

    // this go routine acts as an asynchronous event emitter
	go func() {
		sampleChannel <- "some text"
		close(sampleChannel) // best practice to close the channel by sender
	}()

	wg.Wait()
}
```

<p align="center">
  <img
    src="https://raw.githubusercontent.com/ghoshRitesh12/govents/main/article/meme.jpg"
    alt="omni_man_meme"
  />
</p>

As you notice, it's a lot of Go code to have just a simple event-driven task. I have been wanting to make a Go package, or any package, for instance. So I thought, why not try to implement Node's Event Emitter API, which will provide us with Node like event constructs in Go by using the same code snippet? This package could form a wrapper, which could expose NodeJS like methods for ease of use.

Before trying to implement this package, I Googled for Go packages with similar functionality, and yes, this package is just another one that adds to the heap (pun intended), although I had fun while writing something like this.

## Enter Govents

Govents is a small package that implements the NodeJS like Event Emitter API, which provides Node like event constructs in Go with support for generics. Well, the reason I am explicitly mentioning **generics** is because I have noticed packages with similar implementations that didn't have support for generics.

> If you are curious about the docs then here they are:
>
> - [go doc](https://pkg.go.dev/github.com/ghoshRitesh12/govents)
> - [github repo](https://github.com/ghoshRitesh12/govents)

Here is a basic example to showcase how it works:

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ghoshRitesh12/govents"
)

func main() {
	sample := govents.NewEventEmitter[int]()
	strs := []string{}

	sample.On("data", func(vals ...int) {
		for _, val := range vals {
			strs = append(strs, fmt.Sprintf("v%d", val))
		}
		fmt.Println("within eventListener", strs)
	})

	fmt.Println("outside eventListener", strs)

	for range 4 {
		time.Sleep(time.Second)
		sample.Emit("data", rand.Intn(21), rand.Intn(21), rand.Intn(21))
	}
}
```

The above code prints the following to the console:

```bash
outside eventListener []
within eventListener [v12 v10 v11]
within eventListener [v12 v10 v11 v6 v6 v7]
within eventListener [v12 v10 v11 v6 v6 v7 v8 v18 v16]
```

The `sample.On("data", func(vals ...int) {})` and `sample.Emit("data")` might feel familiar to you if you know Node.

### Event listeners

The `emitter.On(eventName string, cb Listener[T]) error` method takes in two parameters, the first being the `eventName` and the second being the `eventListener` that will be associated with the event. It basically registers an event listener for the event named `eventName` with listener `cb` of type `T`. This method may return an error based on any three conditions:

1. `ErrDuplicateListeners`: We can't register more than one event listener for a single event; if tried, it returns an error.
2. `ErrNoEventName`: This error is returned when the caller passes an empty string as an `eventName` argument.
3. `ErrMaxListenerLimit`: This error is returned when the Event Emitter instance reaches its `maxEventListenerLimit`, which is 10 by default and can be expanded upon by using the `emitter.SetMaxEventListeners(int32) method`.

Well, it wouldn't be like NodeJS's Event Emitter if there wasn't the alias of `emitter.AddEventListener(eventName string, cb Listener[T]) error`, which does exactly what `emitter.On` does.

### Event emitters

The `emitter.Emit(eventName string, vals ...T) error` majorly takes in one parameter being the `eventName` and an optional number of other parameters as `vals` of type `T`. It emits an event named `eventName`, thereby performing an asynchronous call to its listener. This method may also return an error based on any three conditions:

1. `ErrNoEventName`: This error is returned when the caller passes an empty string as an `eventName` argument.
2. `ErrNoEventFound`: This error is returned when the passed in `eventName` doesn't have a registered event listener associated with it.
3. `ErrMaxListenerLimit`: This error is thrown when the Event Emitter instance reaches its `maxEventListenerLimit`.

### Event listeners that run once

Similar to `emitter.On`, there's also an `emitter.Once(eventName string, cb Listener[T]) error`, which takes in the same number of arguments as its counterpart, but this event listener is run only once. If this event is emitted multiple times, then it returns an error of `ErrNoEventFound` on the second emit or call.

### Cleaning event listeners

Node's Event Emitter has `emitter.off(name, callback)` for cleaning up resources and removing the event listener. Govents also has `emitter.Off(eventName string)`, which has the simple job of removing the event listener associated with the `eventName`. In case you didn't notice, it doesn't have the event listener as a second parameter, and it also doesn't check for the event's existence and hence doesn't return any error.

Alias time! Govents also has `emitter.RemoveEventListener(eventName string)` to keep up with Node's `removeEventListener(name, callback)` alias.

### "Burn them all" -The Mad King

> Well, not in a literal sense; I hope there's some Game of Thrones humor in here ;)

Let me clarify the metaphor by mentioning the method `emitter.RemoveAllListeners()`, which can be used to remove all registered event listeners, and yes, just like `emitter.Off()`, it also doesn't return any errors.

### The Utils

I added some utility methods, just in case.

- The `emitter.GetEventNames() []string` gets the names for all the registered events for an Event Emitter instance and returns a slice of strings.

- The `emitter.Len() int32` gets the number of event listeners registered for an Event Emitter instance and returns an int32.

- The `emitter.SetMaxEventListeners(maxListeners int32)` sets the maximum number of event listeners for an Event Emitter instance.

- The `emitter.GetMaxEventListeners() int32` returns the maximum number of event listeners for an Event Emitter instance.

## Wrapping Up

I had tons of fun experimenting with the trial-and-error method of learning things and learned some new things about Go's concurrency model.
This is literally my second Go article, so having some feedback is much appreciated. You could provide feedback on my code, the way I wrote this article, or anything, really :).

\#Go \#Go-Concurrency \#EventEmitter \#NodeJS \#Events
