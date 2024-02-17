package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ghoshRitesh12/govents"
)

func Basic() {
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
