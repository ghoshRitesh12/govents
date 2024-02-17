package govents

import (
	"sync"
	"sync/atomic"
)

type EventEmitter[T comparable] struct {
	// maxListeners(default 10) are always less than 1
	maxEventListeners int32
	ch                chan eventChannel[T]
	emux              *sync.Mutex
	events            map[string]event[T]
}

type (
	Listener[T comparable] func(vals ...T)

	event[T comparable] struct {
		isOnce bool
		cb     Listener[T]
	}

	eventChannel[T comparable] struct {
		eventName string
		cbFnArgs  []T
	}
)

func NewEventEmitter[T comparable]() *EventEmitter[T] {
	eventEmitter := &EventEmitter[T]{
		maxEventListeners: 11,
		emux:              &sync.Mutex{},
		events:            make(map[string]event[T]),
	}

	eventEmitter.ch = make(
		chan eventChannel[T],
	)

	go listenTo[T](eventEmitter)

	return eventEmitter
}

// Registers an event with *eventName* as name and its listener as *cb*.
func (e *EventEmitter[T]) On(eventName string, cb Listener[T]) error {
	e.emux.Lock()
	defer e.emux.Unlock()

	event, err := initEvent[T](e, eventName, cb, false)
	if err != nil {
		return err
	}

	e.events[eventName] = event
	return nil
}

// Registers an event with *eventName* as name and its listener as *cb*, that runs only once.
// Duplicate calls will lead in an event doesn't exist error.
func (e *EventEmitter[T]) Once(eventName string, cb Listener[T]) error {
	e.emux.Lock()
	defer e.emux.Unlock()

	event, err := initEvent[T](e, eventName, cb, true)
	if err != nil {
		return err
	}

	e.events[eventName] = event
	return nil
}

// Emits or calls event named *eventName* and *vals* as arguments.
func (e *EventEmitter[T]) Emit(eventName string, vals ...T) error {
	e.emux.Lock()
	defer e.emux.Unlock()

	if eventName == "" {
		return ErrNoEventName
	}

	_, ok := e.events[eventName]
	if !ok {
		return ErrNoEventFound(eventName)
	}

	if e.Len() >= e.maxEventListeners {
		close(e.ch)
		return ErrMaxListenerLimit
	}

	e.ch <- eventChannel[T]{
		eventName: eventName,
		cbFnArgs:  vals,
	}

	return nil
}

// Removes event listener associated with *eventName*.
func (e *EventEmitter[T]) Off(eventName string) {
	e.emux.Lock()
	defer e.emux.Unlock()

	delete(e.events, eventName)
}

// Registers an event with *eventName* as name and its listener as *cb*.
// Alias for eventEmitter.On() method.
func (e *EventEmitter[T]) AddEventListener(eventName string, cb Listener[T]) error {
	return e.On(eventName, cb)
}

// Removes event listener associated with *eventName*.
// Alias for eventEmitter.Off() method.
func (e *EventEmitter[T]) RemoveEventListener(eventName string) {
	e.Off(eventName)
}

// Removes all registered event listeners for an EventEmitter.
func (e *EventEmitter[T]) RemoveAllListeners() {
	e.emux.Lock()
	defer e.emux.Unlock()

	allEventNames := e.GetEventNames()

	if len(allEventNames) < 1 {
		return
	}

	for _, eventName := range allEventNames {
		delete(e.events, eventName)
	}
}

// Gets number of event listeners registered for an EventEmitter.
func (e *EventEmitter[T]) Len() int32 {
	length := int32(len(e.events))
	return atomic.LoadInt32(&length)
}

// Gets registered event names for an EventEmitter.
func (e *EventEmitter[T]) GetEventNames() []string {
	e.emux.Lock()
	defer e.emux.Unlock()

	allEventNames := []string{}

	if e.Len() < 1 {
		return allEventNames
	}

	for eventName := range e.events {
		allEventNames = append(allEventNames, eventName)
	}

	return allEventNames
}

// Sets maximum number of event listeners for an EventEmitter.
func (e *EventEmitter[T]) SetMaxEventListeners(maxListeners int32) {
	atomic.StoreInt32(&e.maxEventListeners, maxListeners)
}

// Gets maximum number of event listeners for an EventEmitter.
func (e *EventEmitter[T]) GetMaxEventListeners() int32 {
	return atomic.LoadInt32(&e.maxEventListeners)
}
