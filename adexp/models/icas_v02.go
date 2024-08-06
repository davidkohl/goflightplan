package models

import (
	"errors"
	"fmt"

	"gitlab.com/davidkohl/goflightplan/adexp/schema"
)

type IcasV02Model struct {
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

func (data *IcasV02Model) Write(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("empty message")
	}
	//
	for _, item := range schema.IcasV02.Items {
		var err error
		switch item.FRN {
		case 1:
			data.TITLE, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 2:
			data.REFDATA, err = parseREFDATA(s)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 3:
			data.ARCID, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 4:
			data.SSRCODE, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 5:
			data.ADEP, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 6:
			data.ADES, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 7:
			data.ARCTYP, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 8:
			data.IFPLID, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		case 9:
			data.ROUTE, err = parseBasicField(s, item)
			if item.Mendatory && errors.Is(err, schema.ErrorFieldNotPresent) {
				return fmt.Errorf("%v: %s", err, item.DataItem)
			}
		default:
			fmt.Printf("Unrecognized field %v\n", item.DataItem)
		}
	}
	return nil
}
