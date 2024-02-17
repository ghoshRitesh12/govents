package main

import (
	"fmt"

	"github.com/ghoshRitesh12/govents"
)

func Extras() {
	type foo struct {
		bar string
		far int
	}

	process := govents.NewEventEmitter[foo]()

	process.On("enter", func(vals ...foo) {
		for _, f := range vals {
			fmt.Println(f.bar, f.far)
		}
	})

	process.On("exit", func(vals ...foo) {
		fmt.Println("exiting")
	})

	fmt.Println("events registered ->", process.GetEventNames())
	fmt.Println("number of event listeners ->", process.Len())

	fmt.Println("max event listeners before ->", process.GetMaxEventListeners())
	process.SetMaxEventListeners(20)
	fmt.Println("max event listeners after ->", process.GetMaxEventListeners())

	process.Emit("enter", foo{
		bar: "something",
		far: 420,
	})
	process.Emit("exit")
}
