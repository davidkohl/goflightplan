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
	ALDT    string
	ATD     string
	ATOT    string
	CTOT    string `json:"ctot,omitempty"`
	NEWCTOT string `json:"newctot,omitempty"`
	IFPLID  string `json:"ifplid,omitempty"`
	SID     string `json:"sid,omitempty"`
	STAR    string `json:"star,omitempty"`
	REG     string `json:"reg,omitempty"`
	DOF     string
	RMK     string
}

type REFDATA struct {
	SENDER string
	RECVR  string
	SEQNUM string
}

// func parseREFDATA(s string) (REFDATA, error) {

// 	var refdata REFDATA
// 	start := strings.Index(s, "REFDATA")
// 	if start == -1 {
// 		return refdata, fmt.Errorf("refdata not found")
// 	}
// 	iSEQNUM := strings.Index(s[:], "SEQNUM")
// 	evalString := s[start : iSEQNUM+10]
// 	re := regexp.MustCompile(`-\w+\s*-?\w*\s*[A-Z0-9]*`)

// 	// Find all matches
// 	matches := re.FindAllString(evalString, -1)

// 	for _, match := range matches {
// 		parts := regexp.MustCompile(`\s+`).Split(match, -1)
// 		if len(parts) < 2 {
// 			continue
// 		}
// 		switch parts[0] {
// 		case "-SENDER":
// 			if len(parts) > 2 && parts[1] == "-FAC" {
// 				refdata.SENDER = parts[2]
// 			}
// 		case "-RECVR":
// 			if len(parts) > 2 && parts[1] == "-FAC" {
// 				refdata.RECVR = parts[2]
// 			}
// 		case "-SEQNUM":
// 			if len(parts) > 1 {
// 				refdata.SEQNUM = parts[1]
// 			}
// 		}
// 	}

// 	return refdata, nil
// }
