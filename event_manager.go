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
