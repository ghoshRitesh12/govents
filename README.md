<p align="center">
  <img
    src="https://github.com/ghoshRitesh12/govents/assets/101876769/a881784e-681f-4cc5-9bcb-7cd0ab7fae2d"
    alt="montre_go"
  />
</p>

# Govents

A small package that implements NodeJS like events construct in Go, with support for generics.

Go 1.18 or newer version is required; the full documentation is at [https://pkg.go.dev/github.com/ghoshRitesh12/govents](https://pkg.go.dev/github.com/ghoshRitesh12/govents)

## Usage

A basic example:

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

Output for the code above:

```bash
outside eventListener []
within eventListener [v12 v10 v11]
within eventListener [v12 v10 v11 v6 v6 v7]
within eventListener [v12 v10 v11 v6 v6 v7 v8 v18 v16]
```

Some more examples can be found in the [example](https://github.com/ghoshRitesh12/govents/tree/main/example) directory.

Detailed documentation of this package can be found in godoc: [https://pkg.go.dev/github.com/ghoshRitesh12/govents](https://pkg.go.dev/github.com/ghoshRitesh12/govents)
