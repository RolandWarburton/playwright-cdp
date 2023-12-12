package playwright_util

type Emitter struct {
	listeners map[string][]chan string
}

// creates an emitter that can be wired to a channel later
func NewEmitter() *Emitter {
	return &Emitter{
		listeners: make(map[string][]chan string),
	}
}

// register a channel to an emitter
func (e *Emitter) On(event string, listener chan string) {
	e.listeners[event] = append(e.listeners[event], listener)
}

// for any given event, push the data to the listeners
func (e *Emitter) Emit(event string, data string) {
	// for each listener of this event
	// push the data through the channel
	if listeners, ok := e.listeners[event]; ok {
		for _, listener := range listeners {
			go func(listener chan string) {
				listener <- data
			}(listener)
		}
	}
}
