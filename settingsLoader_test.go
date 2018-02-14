package main

import "testing"

func TestNewSettingsLoader(t *testing.T) {
	s := NewSettingsLoader(make(chan Settings), make(chan error), ConfigAWS{})
	if s == nil {
		t.Errorf("could not initialize settings loader")
	}
}

func TestSettingsLoader_Start(t *testing.T) {

}

func TestSettingsLoader_Close(t *testing.T) {

}
