package icao

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"gitlab.com/davidkohl/goflightplan"
)

type ParseHandler struct {
	Fn   func(s string) (*goflightplan.Flightplan, error)
	Name string
}

type ICAOParser struct {
	ParserOpts
	ParseHandlers map[string]ParseHandler
}

type ParserOpts struct {
	CreateMessageSets bool
	AFTNHeader        bool
	ParseHandler      map[string]ParseHandler
}

func NewParser(opts ParserOpts) *ICAOParser {
	p := &ICAOParser{
		ParserOpts:    opts,
		ParseHandlers: make(map[string]ParseHandler),
	}

	p.ParseHandlers["FPL"] = ParseHandler{Name: "FPL", Fn: parseFPL}
	p.ParseHandlers["CHG"] = ParseHandler{Name: "CHG", Fn: parseCHG}
	p.ParseHandlers["CNL"] = ParseHandler{Name: "CNL", Fn: parseCNL}
	p.ParseHandlers["DLA"] = ParseHandler{Name: "DLA", Fn: parseDLA}
	p.ParseHandlers["ARR"] = ParseHandler{Name: "ARR", Fn: parseARR}
	p.ParseHandlers["DEP"] = ParseHandler{Name: "DEP", Fn: parseDEP}

	return p

}

// ParseFPLMessage parses an ICAO FPL message and returns a structured FPLMessage.
func (p *ICAOParser) Parse(s string) (*goflightplan.FlightplanWrapper, error) {
	fplw := &goflightplan.FlightplanWrapper{}

	var fpl *goflightplan.Flightplan = &goflightplan.Flightplan{}
	var send, rec string
	//if AFTNHeader is true, try to extract it
	if p.ParserOpts.AFTNHeader {
		delimStart := "ZCZC"

		//delimEnd := "NNNN"

		i := strings.Index(s, delimStart)
		if i != -1 {
			s := s[i:]
			tmp := strings.Fields(s)
			rec = tmp[3]
			send = tmp[5]
		} else {
			log.Println("AFTNHeader Expected but not found")
		}

	}

	t, err := getTitle(s)
	if err != nil {
		return nil, err
	}

	for _, v := range p.ParseHandlers {
		istart := strings.Index(s, "(")
		iend := strings.Index(s, ")")
		if v.Name == t {
			tmpfpl, err := v.Fn(s[istart : iend+1])
			if err != nil {
				return nil, err
			}
			d, err := json.Marshal(tmpfpl)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(d, fpl)
			if err != nil {
				return nil, err
			}
		}
	}

	if p.ParserOpts.AFTNHeader {
		fpl.REFDATA.SENDER.FAC = send
		fpl.REFDATA.RECVR.FAC = rec
	}

	fplw.Flightplan = *fpl
	return fplw, nil
}

func parseFPL(s string) (*goflightplan.Flightplan, error) {
	var fpl = &goflightplan.Flightplan{}
	start := strings.Index(s, "(FPL-")
	end := strings.Index(s, ")")

	if start == -1 || end == -1 {
		return nil, errors.New("could not get start or end of ICAO message")
	}
	s = s[start:end]
	// Remove the prefix "(FPL-" and the trailing parenthesis ")"
	content := strings.TrimPrefix(s, "(FPL-")
	content = strings.TrimSuffix(content, ")")

	content = strings.ReplaceAll(content, "\n", "")
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
	fpl = &goflightplan.Flightplan{
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
	}

	if err := parseField18(fpl, fields[7]); err != nil {
		return nil, err
	}
	return fpl, nil
}

// ParseCHGMessage parses a CHG message and updates an existing Flightplan.
func parseCHG(message string) (*goflightplan.Flightplan, error) {

	var fp *goflightplan.Flightplan = &goflightplan.Flightplan{}
	if !strings.HasPrefix(message, "(CHG-") {
		return nil, errors.New("message is not a valid CHG message")
	}

	// Remove the prefix and split the message into parts.
	content := strings.TrimPrefix(message, "(CHG-")
	parts := strings.Split(content, "-")

	if len(parts) < 2 {
		return nil, errors.New("incomplete CHG message")
	}

	// The first part after CHG should be the Aircraft ID.
	fp.ARCID = parts[1]

	// Iterate over the remaining parts to handle changes.
	for i := 2; i < len(parts); i++ {
		part := parts[i]
		if pos := strings.Index(part, "/"); pos != -1 {
			fieldNum := part[:pos]
			newValue := part[pos+1:]
			switch fieldNum {
			case "3":
				fp.ARCTYP = newValue
			case "8":
				fp.ADEP = newValue
			case "13":
				fp.ADES = newValue
				if len(newValue) > 4 {
					fp.ADES = newValue[:4]
					fp.EELT = newValue[4:]
				}
			case "15":
				fp.ROUTE = newValue
			case "16":
				fp.ELDT = newValue // Assuming ELDT is relevant here
			case "18":
				if err := parseField18(fp, newValue); err != nil {
				}
			}
		}
	}

	return fp, nil
}
func parseCNL(s string) (*goflightplan.Flightplan, error) {
	var fpl = &goflightplan.Flightplan{}
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	// (CNL-WMT912-EDJA2010-LIRF-DOF/240228)
	fpl.TITLE = "CNL"
	fpl.ARCID = parts[1]
	fpl.ADEP = parts[2][:4]
	fpl.EOBT = parts[2][4:]
	fpl.ADES = parts[3]
	fpl.DOF = strings.Split(parts[4], "/")[1]

	return fpl, nil
}

func parseDLA(s string) (*goflightplan.Flightplan, error) {
	var fpl = &goflightplan.Flightplan{}
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	// (DLA-WZZ5322-LYNI1025-EDJA-DOF/240228)
	fpl.TITLE = "DLA"
	fpl.ARCID = parts[1]
	fpl.ADEP = parts[2][:4]
	fpl.EOBT = parts[2][4:]
	fpl.ADES = parts[3]
	fpl.DOF = strings.Split(parts[4], "/")[1]
	return fpl, nil
}
func parseARR(s string) (*goflightplan.Flightplan, error) {
	var fpl = &goflightplan.Flightplan{}
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	//(ARR-WZZ301-EDJA0910-BKPR1048)
	fpl.TITLE = "ARR"
	fpl.ARCID = parts[1]
	fpl.ADEP = parts[2][:4]
	fpl.EOBT = parts[2][4:]
	fpl.ADES = parts[len(parts)-1][:4]
	fpl.ELDT = parts[len(parts)-1][4:]

	// if len parts > 4, parts[3] is the stop inbetween
	return fpl, nil
}

func parseDEP(s string) (*goflightplan.Flightplan, error) {
	var fpl = &goflightplan.Flightplan{}
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	//(DEP-WZZ456-BKPR1155-EDJA-DOF/240228)
	fpl.TITLE = "DEP"
	fpl.ARCID = parts[1]
	fpl.ADEP = parts[2][:4]
	fpl.EOBT = parts[2][4:]
	fpl.ADES = parts[3][:4]
	fpl.ELDT = parts[3][4:]
	fpl.DOF = strings.Split(parts[4], "/")[1]
	return fpl, nil
}

// ParseOtherInfo handles the parsing of Field 18 and sets the corresponding fields in the Flightplan.
func parseField18(fp *goflightplan.Flightplan, otherInfo string) error {
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

func getTitle(s string) (string, error) {
	// Compile the regex pattern
	pattern := `\(.*?-`
	re := regexp.MustCompile(pattern)

	// Find the first match
	match := re.FindString(s)
	if match == "" {
		return "", fmt.Errorf("no match found")
	}

	// Extract the required part from the match
	// Remove the '(' and '-' from the match
	cleaned := match[1 : len(match)-1] // remove the leading '(' and trailing '-'
	return cleaned, nil
}
