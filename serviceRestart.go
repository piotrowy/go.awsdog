package main

import (
	"fmt"
	"time"
)

type (
	ServiceRestarted struct {
		Name     string
		Attempts int
		Result   string
	}

	ServiceRestartReq struct {
		Name          string
		NumOfAttempts int
		NumOfSecWait  int64
	}

	ServiceRestartHandler struct {
		router   EventRouter
		logCh    chan<- interface{}
		eventBus chan Event
		close    chan bool
	}
)

func (s ServiceRestartReq) String() string {
	return fmt.Sprintf("Request for restarting service %v every %d seconds %d times.", s.Name, s.NumOfSecWait, s.NumOfAttempts)
}

func (s ServiceRestarted) String() string {
	return fmt.Sprintf("Service %v restarted with result %v after %d times.", s.Name, s.Attempts, s.Result)
}

func NewServiceRestartHandler(router EventRouter, logCh chan<- interface{}, bufferSize int) EventHandler {
	return ServiceRestartHandler{
		router:   router,
		logCh:    logCh,
		eventBus: make(chan Event, bufferSize),
		close:    make(chan bool),
	}
}

func (srh ServiceRestartHandler) Start() {
	for {
		select {
		case e := <-srh.eventBus:
			go func() {
				if serviceRestartReq, ok := e.(ServiceRestartReq); ok {
					srh.restart(serviceRestartReq)
				}
			}()
		case <-srh.close:
			return
		}
	}
}

func (srh ServiceRestartHandler) CanHandle(e Event) bool {
	_, ok := e.(ServiceRestartReq)
	return ok
}

func (srh ServiceRestartHandler) GetRoute() chan<- Event {
	return srh.eventBus
}

func (srh ServiceRestartHandler) Close() {
	srh.close <- true
}

func (srh ServiceRestartHandler) restart(event ServiceRestartReq) {
	var (
		i         = 0
		restarted = false
	)
	for range time.Tick(time.Duration(int64(time.Second) * event.NumOfSecWait)) {
		i++
		if srh.try() {
			restarted = true
			break
		}

		if i == event.NumOfAttempts {
			break
		}
	}

	result := ServiceRestarted{
		Name:     event.Name,
		Attempts: i,
	}
	if restarted {
		result.Result = "SUCCESS"
	} else {
		result.Result = "FAILURE"
	}
	srh.logCh <- result
	srh.notify(result)
	srh.router.Route(result)
}

func (srh ServiceRestartHandler) try() bool {
	return true //todo restarting
}

func (srh ServiceRestartHandler) notify(event ServiceRestarted) {
	//todo notifier
}
