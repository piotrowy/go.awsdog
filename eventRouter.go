package main

type (
	Event interface {
		String() string
	}

	EventHandler interface {
		Start()
		GetRoute() chan<- Event
		CanHandle(e Event) bool
		Close()
	}

	EventRouter struct {
		Handlers []EventHandler
	}
)

func NewEventRouter() EventRouter {
	return EventRouter{
		Handlers: []EventHandler{},
	}
}

func (er *EventRouter) Route(e Event) {
	for _, v := range er.Handlers {
		if v.CanHandle(e) {
			v.GetRoute() <- e
		}
	}
}

func (er *EventRouter) RegisterHandler(eh EventHandler) {
	er.Handlers = append(er.Handlers, eh)
}

func (er *EventRouter) Close() {
	for _, v := range er.Handlers {
		v.Close()
	}
}
