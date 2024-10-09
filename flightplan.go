package goflightplan

import (
	"strings"
)

const (
	MessageTypeICAO = iota
	MessageTypeADEXP
)

func GetFlightplanFormat(s string) uint {
	if strings.Contains(s, "-TITLE ") {
		return 1
	}
	return 0
}

type FlightplanWrapper struct {
	Flightplan map[string]interface{}
	Meta       map[string]interface{}
	Raw        string
}

func NewFlightplanWrapper() *FlightplanWrapper {
	meta := make(map[string]interface{}, 0)
	return &FlightplanWrapper{Meta: meta}
}
