package main

import "testing"

func TestEventRouter_Route(t *testing.T) {
	var (
		er = NewEventRouter()
		th = testHandler{}
	)
	er.Route(testEvent{})

	if !th.routed {
		t.Error("Event router did not route event properly.")
	}
}

func TestEventRouter_RegisterHandler(t *testing.T) {
	var (
		er = NewEventRouter()
		th = testHandler{}
	)

	er.RegisterHandler(&th)

	if er.Handlers[0] != &th {
		t.Error("Handler isn't registered.")
	}
}

func TestEventRouter_Close(t *testing.T) {
	var (
		er = NewEventRouter()
		th = testHandler{}
	)
	er.RegisterHandler(&th)

	er.Close()

	if !th.close {
		t.Error("Event router did not close handlers properly.")
	}
}

type (
	testEvent struct{}

	testHandler struct {
		close, routed bool
		route         chan Event
	}
)

func (te testEvent) String() string {
	return ""
}

func (th *testHandler) Start() {
	<-th.route
	th.routed = true
}

func (th *testHandler) GetRoute() chan<- Event {
	return th.route
}

func (th *testHandler) CanHandle(e Event) bool {
	return true
}

func (th *testHandler) Close() {
	th.close = true
}
