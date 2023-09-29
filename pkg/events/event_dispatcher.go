package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")
var ErrHandlerNotFound = errors.New("handler not found")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if has, _ := ed.Has(eventName, handler); has {
		return ErrHandlerAlreadyRegistered
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) (bool, int) {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				return true, i
			}
		}
	}
	return false, -1
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	has, index := ed.Has(eventName, handler)
	if has {
		ed.handlers[eventName] = append(ed.handlers[eventName][:index], ed.handlers[eventName][index+1:]...)
		return nil
	}
	return ErrHandlerNotFound
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandlerInterface)
}
