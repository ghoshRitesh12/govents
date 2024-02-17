package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ghoshRitesh12/govents"
)

func Listeners() {
	sample := govents.NewEventEmitter[string]()
	strs := []string{}

	// same as sample.On("entry", func(vals ...int))
	sample.AddEventListener("entry", func(vals ...string) {
		fmt.Println("within entry event listener")
	})

	sample.On("data", func(vals ...string) {
		strs = append(strs, vals...)
		fmt.Println("within eventListener", strs)
	})

	// this event listener is executed only once
	// multiple emits results in error
	sample.Once("data:once", func(vals ...string) {
		fmt.Println("within eventListener", strs)
	})

	sample.Emit("data:once")        // this emit works fine
	err := sample.Emit("data:once") // this emit causes error
	if err != nil {
		log.Fatalln(err)
	}

	// removes event listener associated with event named "data"
	sample.Off("data")
	// same as above, just an alias for Off()
	sample.RemoveEventListener("entry")

	for range 4 {
		time.Sleep(time.Second)

		// throws error as it tries to emit a non-existent event
		err := sample.Emit("data", randStringInt(), randStringInt(), randStringInt())
		if err != nil {
			log.Fatalln(err)
		}
	}

	// removes all registerd event listeners
	sample.RemoveAllListeners()
}

func randStringInt() string {
	return fmt.Sprintf("%d", rand.Intn(21))
}
