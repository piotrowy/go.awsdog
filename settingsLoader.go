package main

import (
	"net"
	"time"
)

type (
	Settings struct {
		ID            int
		Services      []string
		NumOfSecCheck int64
		NumOfSecWait  int64
		NumOfAttempts int
	}

	SettingsChanged struct {
		Settings
		Date time.Time
	}

	SettingsLoader interface {
		Start()
		Close()
	}

	settingsLoader struct {
		router EventRouter
		errors chan<- interface{}
		close  chan bool
	}
)

func (s Settings) String() string {
	return ""
}

func GetSettings(pk string) Settings {
	return Settings{
		ID:            0,
		Services:      nil,
		NumOfSecCheck: 0,
		NumOfSecWait:  0,
		NumOfAttempts: 0,
	}
}

func NewSettingsLoader(router EventRouter, errors chan<- interface{}, config ConfigAWS) SettingsLoader {
	return settingsLoader{
		router: router,
		errors: errors,
		close:  make(chan bool),
	}
}

func (s settingsLoader) Start() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		select {
		case <-s.close:
			return
		default:
			if conn, err := ln.Accept(); err != nil {
				s.errors <- err
			} else {
				s.readData(conn)
			}
		}
	}
}

func (s settingsLoader) Close() {
	s.close <- true
}

func (s settingsLoader) readData(conn net.Conn) {
	buf := make([]byte, 32)
	if n, err := conn.Read(buf); err != nil {
		s.errors <- err
	} else {
		if string(buf[:n]) == "new settings" {
			go s.fetchDynamoDB()
		}
	}
}

func (s settingsLoader) fetchDynamoDB() {
	//todo fetcher
	s.router.Route(SettingsChanged{
		Settings: Settings{
			ID:            0,
			Services:      nil,
			NumOfSecCheck: 0,
			NumOfSecWait:  0,
			NumOfAttempts: 0,
		},
		Date: time.Time{},
	})
}
