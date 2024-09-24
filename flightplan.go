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

type CFL struct {
	FL string `json:"fl,omitempty"`
}

type FlightplanWrapper struct {
	Flightplan Flightplan
	Meta       map[string]interface{}
	Raw        string
}

func NewFlightplanWrapper() *FlightplanWrapper {
	meta := make(map[string]interface{}, 0)
	return &FlightplanWrapper{Meta: meta}
}

type Flightplan struct {
	TITLE   string  `json:"title,omitempty"`
	REFDATA REFDATA `json:"refdata,omitempty"`
	ARCID   string  `json:"arcid,omitempty"`
	SSRCODE string  `json:"ssrcode,omitempty"`
	ADEP    string  `json:"adep,omitempty"`
	ADES    string  `json:"ades,omitempty"`
	ARCTYP  string  `json:"arctyp,omitempty"`
	WKTRC   string  `json:"wktrc,omitempty"`
	ROUTE   string  `json:"route,omitempty"`
	FLTRUL  string  `json:"fltrul,omitempty"`
	FLTTYP  string  `json:"flttyp,omitempty"`
	FPLCAT  string  `json:"fplcat,omitempty"`
	RFL     string  `json:"rfl,omitempty"`
	CFL     CFL     `json:"cfl,omitempty"`
	ICA     string  `json:"ica,omitempty"`
	EOBT    string  `json:"eobt,omitempty"`
	EOBD    string  `json:"eobd,omitempty"`
	EELT    string  `json:"eelt,omitempty"`
	ELDT    string  `json:"eldt,omitempty"`
	ALDT    string  `json:"aldt,omitempty"`
	ATD     string  `json:"atd,omitempty"`
	ATOT    string  `json:"atot,omitempty"`
	CTOT    string  `json:"ctot,omitempty"`
	NEWCTOT string  `json:"newctot,omitempty"`
	IFPLID  string  `json:"ifplid,omitempty"`
	SID     string  `json:"sid,omitempty"`
	STAR    string  `json:"star,omitempty"`
	REG     string  `json:"reg,omitempty"`
	DOF     string  `json:"dof,omitempty"`
	RMK     string  `json:"rmk,omitempty"`
}

type REFDATA struct {
	SENDER string `json:"sender,omitempty"`
	RECVR  string `json:"recvr,omitempty"`
	SEQNUM string `json:"seqnum,omitempty"`
}
