package goflightplan

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"gitlab.com/davidkohl/goflightplan/adexp"
	"gitlab.com/davidkohl/goflightplan/adexp/icasv01"
)

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
	CFL     string  `json:"cfl,omitempty"`
	ICA     string  `json:"ica,omitempty"`
	EOBT    string  `json:"eobt,omitempty"`
	EOBD    string  `json:"eobd,omitempty"`
	ELDT    string  `json:"eldt,omitempty"`
	CTOT    string  `json:"ctot,omitempty"`
	NEWCTOT string  `json:"newctot,omitempty"`
	IFPLID  string  `json:"ifplid,omitempty"`
	SID     string  `json:"sid,omitempty"`
	STAR    string  `json:"star,omitempty"`
	REG     string  `json:"reg,omitempty"`
}

type REFDATA struct {
	SENDER string
	RECVR  string
	SEQNUM string
}

func ParseREFDATA(s string) (REFDATA, error) {

	var refdata REFDATA
	start := strings.Index(s, "REFDATA")
	if start == -1 {
		return refdata, fmt.Errorf("refdata not found")
	}
	iSEQNUM := strings.Index(s[:], "SEQNUM")
	evalString := s[start : iSEQNUM+10]
	re := regexp.MustCompile(`-\w+\s*-?\w*\s*[A-Z0-9]*`)

	// Find all matches
	matches := re.FindAllString(evalString, -1)

	for _, match := range matches {
		parts := regexp.MustCompile(`\s+`).Split(match, -1)
		if len(parts) < 2 {
			continue
		}
		switch parts[0] {
		case "-SENDER":
			if len(parts) > 2 && parts[1] == "-FAC" {
				refdata.SENDER = parts[2]
			}
		case "-RECVR":
			if len(parts) > 2 && parts[1] == "-FAC" {
				refdata.RECVR = parts[2]
			}
		case "-SEQNUM":
			if len(parts) > 1 {
				refdata.SEQNUM = parts[1]
			}
		}
	}

	return refdata, nil
}

func parseBasicField(s string, f adexp.DataField) (string, error) {
	start := strings.Index(s, fmt.Sprintf("-%s ", f.DataItem))
	if start == -1 {
		return "", adexp.ErrorFieldNotPresent
	}
	s = s[start:]
	next := strings.Index(s[1:], "-")
	if next == -1 {
		next = len(s)
	}
	sub := strings.SplitN(s[:next], " ", 2)
	return strings.TrimSpace(sub[1]), nil
}

func (fpl *Flightplan) Write(s string) error {
	//check if string is valid
	if len(s) == 0 {
		return fmt.Errorf("could not parse flightplan string. message string is empty")
	}
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "NNNN")
	s = strings.TrimSpace(s)
	// get format of flightplan: ICAO or ADEXP
	t := ""
	if strings.Contains(s, "-TITLE") {
		t = "ADEXP"
		_ = t
		d := reflect.ValueOf(fpl)
		e := d.Elem()

		for _, item := range icasv01.BFD.Items {
			f := e.FieldByName(item.DataItem)
			if !f.IsValid() {
				log.Printf("Skipped %s bacause is not valid field", item.DataItem)
				continue
			}
			if !f.CanSet() {
				log.Printf("Skipped %s bacause can not be set", item.DataItem)
				continue
			}

			switch item.Type {
			case 0:
				g, err := parseBasicField(s, item)
				if errors.Is(err, adexp.ErrorFieldNotPresent) && item.Mendatory {
					return fmt.Errorf("error: mendatory %s: %s", err, item.DataItem)
				}
				f.SetString(g)
			case 1:
				//log.Printf("List field %s needs to be implemented", item.DataItem)
			case 2:
				//log.Printf("Structured field %s needs to be implemented", item.DataItem)
			default:
			}
		}
		return nil
	}

	// select message set to chose from

	//get identifier (BFD,CFD .. or FPL,CHG,)

	return nil
}
