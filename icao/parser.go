package icao

import (
	"errors"
	"reflect"
	"strings"

	"gitlab.com/davidkohl/goflightplan"
)

type ICAOParser struct {
	ParserOpts
}

type ParserOpts struct {
	CreateMessageSets bool
	AFTNHeader        bool
}

func NewParser(opts ParserOpts) *ICAOParser {
	return &ICAOParser{
		ParserOpts: opts,
	}

}

// ParseFPLMessage parses an ICAO FPL message and returns a structured FPLMessage.
func (p *ICAOParser) Parse(s string) (*goflightplan.FlightplanWrapper, error) {
	start := strings.Index(s, "(FPL-")
	end := strings.Index(s, ")")
	if start == -1 || end == -1 {
		return nil, errors.New("could not get start or end of ICAO message")
	}
	s = s[start:end]
	// Remove the prefix "(FPL-" and the trailing parenthesis ")"
	content := strings.TrimPrefix(s, "(FPL-")
	content = strings.TrimSuffix(content, ")")

	// Split the fields using "-" as the delimiter
	fields := strings.Split(content, "-")

	if len(fields) < 7 {
		return nil, errors.New("incomplete FPL message")
	}

	// Split the aircraft type and wake turbulence category
	aircraftParts := strings.Split(fields[2], "/")
	if len(aircraftParts) != 2 {
		return nil, errors.New("invalid aircraft type or wake turbulence category format")
	}

	// Extract the departure aerodrome and estimated off-block time (EOBT)
	departureInfo := fields[4]

	// Extract route, destination aerodrome, and estimated elapsed time
	route := fields[5]
	destinationInfo := fields[6]
	if len(destinationInfo) < 7 {
		return nil, errors.New("destination information is too short")
	}

	// Ensure EELT is assigned correctly if present
	eelt := ""
	_ = eelt
	if len(destinationInfo) > 4 {
		eelt = destinationInfo[4:]
	}

	// Extract and assign values to the FPLMessage struct
	fpl := &goflightplan.Flightplan{
		TITLE:  "FPL",
		ARCID:  fields[0],            // Aircraft ID
		FLTRUL: string(fields[1][0]), // Flight rules
		FLTTYP: string(fields[1][1]), // Flight type
		ARCTYP: aircraftParts[0],     // Aircraft type
		WKTRC:  aircraftParts[1],     // Wake turbulence category
		ADEP:   departureInfo[:4],    // Departure aerodrome
		EOBT:   departureInfo[4:],    // Estimated Off-Block Time (departure time)
		ROUTE:  route,                // Route
		ADES:   destinationInfo[:4],  // Destination aerodrome
		EELT:   eelt,                 // Estimated Elapsed Time
		//OtherInformation: fields[6],            // Other information
	}

	if err := ParseOtherInfo(fpl, fields[7]); err != nil {
		return nil, err
	}
	fplw := goflightplan.NewFlightplanWrapper()
	fplw.Flightplan = *fpl
	return fplw, nil
}

// ParseOtherInfo handles the parsing of Field 18 and sets the corresponding fields in the Flightplan.
func ParseOtherInfo(fp *goflightplan.Flightplan, otherInfo string) error {
	infoMap := make(map[string]string)
	var currentKey string
	var currentVal strings.Builder

	tokens := strings.Fields(otherInfo)
	for _, token := range tokens {
		if pos := strings.Index(token, "/"); pos != -1 {
			if currentKey != "" {
				infoMap[currentKey] = strings.TrimSpace(currentVal.String())
				currentVal.Reset()
			}
			currentKey = token[:pos]
			currentVal.WriteString(token[pos+1:] + " ")
		} else {
			currentVal.WriteString(token + " ")
		}
	}

	if currentKey != "" {
		infoMap[currentKey] = strings.TrimSpace(currentVal.String())
	}

	val := reflect.ValueOf(fp).Elem()
	for k, v := range infoMap {
		if fieldVal := val.FieldByName(strings.ToUpper(k)); fieldVal.IsValid() && fieldVal.CanSet() {
			if fieldVal.Kind() == reflect.String {
				fieldVal.SetString(v)
			}
		}
	}

	return nil
}
