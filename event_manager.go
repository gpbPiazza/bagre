package main

type Listener interface {
	Subscribe(et EventType, payload any)
	ID() int
}

type EventType int

const (
	animationEnded EventType = iota
	removeUnit
)

type EventManager struct {
	listByEventType map[EventType][]Listener
}

func NewEventManager() *EventManager {
	return &EventManager{
		listByEventType: make(map[EventType][]Listener),
	}
}

func (a *EventManager) PublishByID(et EventType, payload any, id int) {
	panic("not implemented")
}

func (a *EventManager) Publish(et EventType, payload any) {
	for _, l := range a.listByEventType[et] {
		l.Subscribe(et, payload)
	}
}

func (a *EventManager) subscribe(eventType EventType, l Listener) {
	a.listByEventType[eventType] = append(a.listByEventType[eventType], l)
}

func (a *EventManager) unsubscribe(eventType EventType, list Listener) {
	listeners := a.listByEventType[eventType]
	lenList := len(listeners)
	lastIndex := lenList - 1

	for i, l := range listeners {
		if l.ID() == list.ID() {
			listeners[lastIndex], listeners[i] = listeners[i], listeners[lastIndex]
			a.listByEventType[eventType] = listeners[:lastIndex]
		}
	}
}
