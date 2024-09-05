package adexp

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gitlab.com/davidkohl/goflightplan"
)

type ADEXPParser struct {
	MessageSchemaSet []MessageSet
	ParserOpts
}

type ParserOpts struct {
	AFTNHeader bool
}

func NewParser(m []MessageSet, opts ParserOpts) *ADEXPParser {
	var p ADEXPParser = ADEXPParser{
		MessageSchemaSet: m,
		ParserOpts:       opts,
	}

	return &p
}

func getTitle(s string) (string, error) {
	if strings.Contains(s, "-TITLE ") {
		i := strings.Index(s, "-TITLE ")
		start := i + 7
		end := strings.Index(s[start:], "-")
		return s[start : start+end-1], nil
	}
	// Implement ICAO later
	return "", errors.New("could not get message type; TITLE field missing")
}

/*
Takes a raw adexp flight plan message and attemps to parse it.

Will return a flightplanwrapper with the flightplan and additional meta info.
*/
func (p *ADEXPParser) Parse(s string) (*goflightplan.FlightplanWrapper, error) {

	fpl := &goflightplan.Flightplan{}
	var fplwrapper = goflightplan.NewFlightplanWrapper()
	t, err := getTitle(s)
	if err != nil {
		return nil, fmt.Errorf("could not get message type. TITLE field is missing: %s ", s)
	}
	for _, v := range p.MessageSchemaSet {
		_ = v
		_ = t
		for _, v := range v.Set {
			if v.Category != t {
				continue
			}
			// Matching Schema has been found; set meta and start processing.
			fplwrapper.Meta["Setname"] = v.Name
			fplwrapper.Meta["SetCategory"] = v.Category
			fplwrapper.Meta["SetVersion"] = v.Version
			d := reflect.ValueOf(fpl)
			e := d.Elem()
			for _, v := range v.Items {
				var f reflect.Value
				f = e.FieldByName(v.DataItem)
				if v.Target != "" {
					f = e.FieldByName(v.Target)
				}
				if !f.IsValid() {
					log.Printf("Skipped %s bacause is not valid field", f)
					continue
				}
				if !f.CanSet() {
					log.Printf("Skipped %s bacause can not be set", f)
					continue
				}
				e.Type()
				switch v.Type {
				case Basicfield:
					a, err := parseBasicField(s, v)
					if err != nil {
						continue
					}
					f.SetString(a)
				case StructuredField:
					//fmt.Println("NOT YET IMPLEMENTED")

				}
			}
			fplwrapper.Flightplan = *fpl
			fplwrapper.Raw = s
			return fplwrapper, nil
		}
	}

	return nil, fmt.Errorf("no matching schema was found for TITLE: %s", t)
}

func parseBasicField(s string, f DataField) (string, error) {
	start := strings.Index(s, fmt.Sprintf("-%s ", f.DataItem))
	if start == -1 {
		return "", ErrorFieldNotPresent
	}
	s = s[start:]
	next := strings.Index(s[1:], "-")
	if next == -1 {
		next = len(s)
	}
	sub := strings.SplitN(s[:next], " ", 2)
	return strings.TrimSpace(sub[1]), nil
}
