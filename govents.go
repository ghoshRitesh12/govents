package govents

import (
	"errors"
	"fmt"
	"sync"
)

type (
	EventEmitter[T comparable] interface {
		On(string, func(...T)) error
		Emit(string, ...T) error
		Off(string) error              // TODO
		Once(string, func(...T)) error // TODO
		RemoveAllListeners() error     // TODO

		// Listeners(eventName string) []listeners
		GetEventNames() []string // TODO
		Len() int                // TODO

		GetMaxListeners() int      // TODO
		SetMaxListeners(int) error // TODO
	}

	EventChannel[T comparable] struct {
		eventName string
		cbFnArgs  []T
	}

	Event[T comparable] struct {
		noOfEvents int
		maxEvents  int
		_ch        chan EventChannel[T]
		eventMut   *sync.Mutex
		events     map[string]func(vals ...T)
	}
)

func NewEventEmitter[T comparable]() *Event[T] {
	eventEmitter := &Event[T]{
		noOfEvents: 0,
		maxEvents:  10,
		eventMut:   &sync.Mutex{},
		events:     make(map[string]func(vals ...T)),
	}

	eventEmitter._ch = make(
		chan EventChannel[T],
		eventEmitter.maxEvents,
	)

	go func() {
		for e := range eventEmitter._ch {
			eventCb, ok := eventEmitter.events[e.eventName]
			if !ok {
				// fmt.Errorf("no event with name %v is registered", e.eventName)
				continue
			}

			eventCb(e.cbFnArgs...)
		}
	}()

	return eventEmitter
}

func (e *Event[T]) On(eventName string, callbackFn func(vals ...T)) error {
	if eventName == "" {
		return errors.New("no event name found when registering event listener")
	}

	e.eventMut.Lock()
	defer e.eventMut.Unlock()

	if e.noOfEvents == e.maxEvents {
		close(e._ch)
		return errors.New("max number of events reached")
	}

	e.noOfEvents += 1
	e.events[eventName] = callbackFn

	return nil
}

func (e *Event[T]) Emit(eventName string, vals ...T) error {
	if eventName == "" {
		return errors.New("no event name found when emitting event")
	}

	e.eventMut.Lock()
	defer e.eventMut.Unlock()

	_, ok := e.events[eventName]
	if !ok {
		return fmt.Errorf("no event with name %v is registered", eventName)
	}

	e._ch <- EventChannel[T]{
		eventName: eventName,
		cbFnArgs:  vals,
	}

	return nil
}
