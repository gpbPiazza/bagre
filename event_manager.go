package main

type EventHandler interface {
	Handle(et EventType, payload any)
}

type EventType int

const (
	animationEnded EventType = iota
	removeUnit
	attackAnimationEnded
)

type EventManager struct {
	listByEventType map[EventType][]EventHandler
}

func NewEventManager() *EventManager {
	return &EventManager{
		listByEventType: make(map[EventType][]EventHandler),
	}
}

func (a *EventManager) PublishByID(et EventType, payload any, id int) {
	panic("not implemented")
}

func (a *EventManager) Publish(et EventType, payload any) {
	for _, l := range a.listByEventType[et] {
		l.Handle(et, payload)
	}
}

func (a *EventManager) subscribe(eventType EventType, l EventHandler) {
	a.listByEventType[eventType] = append(a.listByEventType[eventType], l)
}
