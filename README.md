<p align="center">
  <img
    src="https://github.com/ghoshRitesh12/govents/assets/101876769/a881784e-681f-4cc5-9bcb-7cd0ab7fae2d"
    alt="govents"
  />
</p>

<h1 align="center"> Govents </h1>

<p align="center">  
	<span>
		<a href="https://pkg.go.dev/github.com/ghoshRitesh12/govents">
			<img
				src="https://pkg.go.dev/badge/github.com/ghoshRitesh12/govents.svg"
				alt="Go Reference"
			/>
		</a>
	</span>
	&nbsp;
	<span>
		<a href="https://riteshghosh.hashnode.dev/govents-nodejs-like-event-emitter-for-go">
			<img 
				alt="Dynamic XML Badge" 
				src="https://img.shields.io/badge/dynamic/xml?url=https%3A%2F%2Friteshghosh.hashnode.dev%2Frss.xml&query=%2F%2Frss%2Fchannel%2Fitem%5B1%5D%2Ftitle&style=social&logoColor=%2389EBFF&label=Article&color=%2389EBFF&link=https%3A%2F%2Friteshghosh.hashnode.dev%2Fgovents-nodejs-like-event-emitter-for-go"
			>
		</a>
	</span>
	&nbsp;
	<span>
		<img
			src="https://img.shields.io/github/v/tag/ghoshRitesh12/govents?style=social"
			alt="Go Reference"
		/>
	</span>
</p>

A small package that implements the NodeJS like Event Emitter API, which provides Node like event constructs in Go with support for generics.

Go 1.18 or newer version is required; the full documentation is at [https://pkg.go.dev/github.com/ghoshRitesh12/govents](https://pkg.go.dev/github.com/ghoshRitesh12/govents)

> I think I wrote a [banger of an article](https://riteshghosh.hashnode.dev/govents-nodejs-like-event-emitter-for-go) up at Hashnode; go check it out if interested.

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
