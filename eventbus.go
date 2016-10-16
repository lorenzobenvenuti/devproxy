package devproxy

import "reflect"

type EventBus interface {
	Subscribe(topic string, handler func(message interface{}))
	Unsubscribe(topic string, handler func(message interface{}))
	Dispatch(topic string, message interface{})
}

type eventBusImpl struct {
	handlers map[string][]func(message interface{})
}

func NewEventBus() EventBus {
	return &eventBusImpl{handlers: make(map[string][]func(message interface{}))}
}

func (eb *eventBusImpl) Subscribe(topic string, handler func(message interface{})) {
	if _, ok := eb.handlers[topic]; ok {
		eb.handlers[topic] = append(eb.handlers[topic], handler)
	} else {
		eb.handlers[topic] = []func(message interface{}){handler}
	}
}

func (eb *eventBusImpl) Unsubscribe(topic string, handler func(message interface{})) {
	if value, ok := eb.handlers[topic]; ok {
		index := -1
		for i, v := range value {
			// Storing handler pointers and comparing them with &handler doesn't work
			// Are functions passed by value? If so, why this works?
			stored := reflect.ValueOf(v)
			current := reflect.ValueOf(handler)
			if stored.Pointer() == current.Pointer() {
				index = i
				break
			}
		}
		if index != -1 {
			eb.handlers[topic] = append(value[:index], value[index+1:]...)
		}
	}
}

func (eb *eventBusImpl) Dispatch(topic string, message interface{}) {
	if value, ok := eb.handlers[topic]; ok {
		for _, v := range value {
			v(message)
		}
	}
}
