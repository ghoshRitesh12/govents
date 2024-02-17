package govents

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateListeners = errors.New("can't have duplicate listeners on one event")
	ErrNoEventName        = errors.New("no event name found when registering event listener")
	ErrMaxListenerLimit   = errors.New("max number of events reached")
	ErrEventEmitterClosed = errors.New("event emitter is closed")
)

func ErrNoEventFound(eventName string) error {
	return fmt.Errorf("no event named %v exists", eventName)
}
