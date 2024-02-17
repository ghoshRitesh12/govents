package govents_test

import (
	"fmt"
	"testing"

	"github.com/ghoshRitesh12/govents"
)

func TestNewEventEmitter(t *testing.T) {
	p := govents.NewEventEmitter[int]()

	err := p.On("data", func(vals ...int) {
		fmt.Println("data vals ->>", vals)
	})
	if err != nil {
		t.Error(err)
	}

	for i := range 12 {
		fmt.Println(i)

		err := p.Emit("data", 69, i)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
