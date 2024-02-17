package govents

import (
	"log"
	"strings"
)

func listenTo[T comparable](eventEmitter *EventEmitter[T]) {
	for {
		func() {
			// eventEmitter.eventsMu.Lock()
			// defer eventEmitter.eventsMu.Unlock()

			e, okay := <-eventEmitter.ch
			if !okay {
				log.Fatalln(ErrEventEmitterClosed)
			}

			event, ok := eventEmitter.events[e.eventName]
			if !ok {
				log.Fatalln(ErrNoEventFound(e.eventName))
			}

			if event.isOnce {
				event.cb(e.cbFnArgs...)
				delete(eventEmitter.events, e.eventName)
				return
			}

			event.cb(e.cbFnArgs...)
		}()
	}
}

// to check for defaults and initialize an event
func initEvent[T comparable](e *EventEmitter[T], eventName string, cb Listener[T], isOnce bool) (event[T], error) {
	eventName = strings.TrimSpace(eventName)
	tempEvent := event[T]{
		isOnce: isOnce,
		cb:     cb,
	}

	if _, eventExist := e.events[eventName]; eventExist {
		return tempEvent, ErrDuplicateListeners
	}

	if eventName == "" {
		return tempEvent, ErrNoEventName
	}

	if e.Len() >= e.maxEventListeners {
		close(e.ch)
		return tempEvent, ErrMaxListenerLimit
	}

	return tempEvent, nil
}
