package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/davidkohl/goflightplan/adexp/schema"
)

type REFDATA struct {
	RECVR  string `json:"recvr,omitempty"`
	SENDER string `json:"sender,omitempty"`
	SEQNUM string `json:"seqnum,omitempty"`
}

func parseREFDATA(s string) (REFDATA, error) {
	s = strings.ReplaceAll(s, "\n", "")
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

func parseBasicField(s string, f string) (string, error) {
	s = strings.ReplaceAll(s, "\n", "")
	start := strings.Index(s, f)
	if start == -1 {
		return "", schema.ErrorFieldNotPresent
	}
	s = s[start:]
	next := strings.Index(s, "-")
	if next == -1 {
		return "", schema.ErrorFieldNotPresent
	}
	sub := strings.SplitN(s[:next], " ", 2)
	return strings.TrimSpace(sub[1]), nil
}

type IcasV01Model struct {
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
	RFL     string  `json:"rfl,omitempty"`
	CFL     string  `json:"cfl,omitempty"`
	ICA     string  `json:"ica,omitempty"`
	EOBT    string  `json:"eobt,omitempty"`
	ELDT    string  `json:"eldt,omitempty"`
	CTOT    string  `json:"ctot,omitempty"`
	IFPLID  string  `json:"ifplid,omitempty"`
}

func (data *IcasV01Model) Write(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("empty message")
	}
	//
	for _, item := range schema.IcasV01.Items {
		var err error
		switch item.FRN {
		case 1:
			data.TITLE, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, "TITLE")
			}
		case 2:
			data.REFDATA, err = parseREFDATA(s)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		case 3:
			data.ARCID, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		case 4:
			data.SSRCODE, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		case 5:
			data.ADEP, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		case 6:
			data.ADES, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		case 7:
			data.ARCTYP, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		case 8:
			data.IFPLID, err = parseBasicField(s, item.DataItem)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return err
			}
		default:
			//fmt.Printf("Unrecognized field %v\n", item.DataItem)
		}
	}
	return nil
}
