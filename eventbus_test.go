package devproxy

import "testing"

type listener struct {
	foo int
	bar int
}

func (l *listener) onFooMessage(message interface{}) {
	l.foo++
}

func (l *listener) onFooMessageAlt(message interface{}) {
	l.foo++
}

func (l *listener) onBarMessage(message interface{}) {
	l.bar++
}

func TestEventBus(t *testing.T) {
	eventBus := NewEventBus()
	listener := &listener{0, 0}
	eventBus.Subscribe("foo", listener.onFooMessage)
	eventBus.Subscribe("foo", listener.onFooMessageAlt)
	eventBus.Subscribe("bar", listener.onBarMessage)
	eventBus.Dispatch("foo", "A Foo message")
	eventBus.Dispatch("bar", "A non foo message")
	eventBus.Unsubscribe("foo", listener.onFooMessage)
	eventBus.Unsubscribe("foo", listener.onFooMessageAlt)
	eventBus.Dispatch("foo", "Another foo message")
	if listener.foo != 2 {
		t.Errorf("Foo listener invoked %d times instead of 1", listener.foo)
	}
	if listener.bar != 1 {
		t.Errorf("Bar  listener invoked %d times instead of 1", listener.bar)
	}
}
