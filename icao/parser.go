package icao

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type ParseHandler struct {
	Fn   func(s string) (map[string]interface{}, error)
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
func (p *ICAOParser) Parse(s string) (map[string]interface{}, error) {
	var fpl map[string]interface{} = make(map[string]interface{})
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
	fmt.Println("TITLE:", t)
	start := strings.Index(s, "(")
	end := strings.Index(s, ")")
	fmt.Println("START", start, "END:", end)
	if start == -1 || end == -1 {
		return nil, errors.New("could not get start or end of ICAO message")
	}
	s = s[start : end+1]
	// Remove the prefix "(FPL-" and the trailing parenthesis ")"
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")

	s = strings.ReplaceAll(s, "\n", "")
	for _, v := range p.ParseHandlers {
		if end == -1 {
			return nil, errors.New("malformed message: missing )")
		}
		if v.Name == t {
			tmpfpl, err := v.Fn(s)
			if err != nil {
				return nil, err
			}
			d, err := json.Marshal(tmpfpl)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(d, &fpl)
			if err != nil {
				return nil, err
			}
		}
	}

	if p.ParserOpts.AFTNHeader {
		fpl["REFDATA"].(map[string]interface{})["SENDER"].(map[string]interface{})["FAC"] = send
		fpl["REFDATA"].(map[string]interface{})["RECVR"].(map[string]interface{})["FAC"] = rec
	}

	return fpl, nil
}

func parseFPL(s string) (map[string]interface{}, error) {
	var fpl = make(map[string]interface{}, 0)

	// Split the fields using "-" as the delimiter
	fields := strings.Split(s, "-")

	if len(fields) < 7 {
		return nil, errors.New("incomplete FPL message")
	}
	fmt.Println(fields)

	// Split the aircraft type and wake turbulence category
	aircraftParts := strings.Split(fields[3], "/")
	if len(aircraftParts) != 2 {
		return nil, errors.New("invalid aircraft type or wake turbulence category format")
	}

	// Extract the departure aerodrome and estimated off-block time (EOBT)
	departureInfo := fields[5]

	// Extract route, destination aerodrome, and estimated elapsed time
	route := fields[6]
	destinationInfo := fields[7]
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

	fpl["TITLE"] = "FPL"
	fpl["ARCID"] = fields[0]             // Aircraft ID
	fpl["FLTRUL"] = string(fields[1][0]) // Flight rules
	fpl["FLTTYP"] = string(fields[1][1]) // Flight type
	fpl["ARCTYP"] = aircraftParts[0]     // Aircraft type
	fpl["WKTRC"] = aircraftParts[1]      // Wake turbulence category
	fpl["ADEP"] = departureInfo[:4]      // Departure aerodrome
	fpl["EOBT"] = departureInfo[4:]      // Estimated Off-Block Time (departure time)
	fpl["ROUTE"] = route                 // Route
	fpl["ADES"] = destinationInfo[:4]    // Destination aerodrome
	fpl["EELT"] = eelt                   // Estimated Elapsed Time

	if err := parseField18(&fpl, fields[8]); err != nil {
		return nil, err
	}
	return fpl, nil
}

// ParseCHGMessage parses a CHG message and updates an existing Flightplan.
func parseCHG(s string) (map[string]interface{}, error) {

	var fpl map[string]interface{} = make(map[string]interface{})

	// Remove the prefix and split the message into parts.

	parts := strings.Split(s, "-")

	if len(parts) < 2 {
		return nil, errors.New("incomplete CHG message")
	}

	// The first part after CHG should be the Aircraft ID.]\
	fpl["TITLE"] = "CHG"
	fpl["ARCID"] = parts[1]

	// Iterate over the remaining parts to handle changes.
	for i := 2; i < len(parts); i++ {
		part := parts[i]
		if pos := strings.Index(part, "/"); pos != -1 {
			fieldNum := part[:pos]
			newValue := part[pos+1:]
			switch fieldNum {
			case "3":
				fpl["ARCTYP"] = newValue
			case "8":
				fpl["ADEP"] = newValue
			case "13":
				fpl["ADES"] = newValue
				if len(newValue) > 4 {
					fpl["ADES"] = newValue[:4]
					fpl["EELT"] = newValue[4:]
				}
			case "15":
				fpl["ROUTE"] = newValue
			case "16":
				fpl["ELDT"] = newValue // Assuming ELDT is relevant here
			case "18":
				if err := parseField18(&fpl, newValue); err != nil {
				}
			}
		}
	}

	return fpl, nil
}
func parseCNL(s string) (map[string]interface{}, error) {
	var fpl = make(map[string]interface{})
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	// (CNL-WMT912-EDJA2010-LIRF-DOF/240228)
	fpl["TITLE"] = "CNL"
	fpl["ARCID"] = parts[1]
	fpl["ADEP"] = parts[2][:4]
	fpl["EOBT"] = parts[2][4:]
	fpl["ADES"] = parts[3]
	fpl["DOF"] = strings.Split(parts[4], "/")[1]

	return fpl, nil
}

func parseDLA(s string) (map[string]interface{}, error) {
	var fpl = make(map[string]interface{})
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	// (DLA-WZZ5322-LYNI1025-EDJA-DOF/240228)
	fpl["TITLE"] = "DLA"
	fpl["ARCID"] = parts[1]
	fpl["ADEP"] = parts[2][:4]
	fpl["EOBT"] = parts[2][4:]
	fpl["ADES"] = parts[3]
	fpl["DOF"] = strings.Split(parts[4], "/")[1]
	return fpl, nil
}
func parseARR(s string) (map[string]interface{}, error) {
	var fpl = make(map[string]interface{})
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	//(ARR-WZZ301-EDJA0910-BKPR1048)
	fpl["TITLE"] = "ARR"
	fpl["ARCID"] = parts[1]
	fpl["ARCID"] = parts[2][:4]
	fpl["EOBT"] = parts[2][4:]
	fpl["ADES"] = parts[len(parts)-1][:4]
	fpl["ELDT"] = parts[len(parts)-1][4:]

	// if len parts > 4, parts[3] is the stop inbetween
	return fpl, nil
}

func parseDEP(s string) (map[string]interface{}, error) {
	var fpl = make(map[string]interface{})
	s = strings.TrimPrefix(s, "(CHG-")
	parts := strings.Split(s, "-")
	//(DEP-WZZ456-BKPR1155-EDJA-DOF/240228)
	fpl["TITLE"] = "DEP"
	fpl["ARCID"] = parts[1]
	fpl["ARCID"] = parts[2][:4]
	fpl["EOBT"] = parts[2][4:]
	fpl["ADES"] = parts[3][:4]
	fpl["EDLT"] = parts[3][4:]
	fpl["DOF"] = strings.Split(parts[4], "/")[1]
	return fpl, nil
}

// ParseOtherInfo handles the parsing of Field 18 and sets the corresponding fields in the Flightplan.
func parseField18(fp *map[string]interface{}, otherInfo string) error {

	var currentKey string
	var currentVal strings.Builder

	tokens := strings.Fields(otherInfo)
	for _, token := range tokens {
		if pos := strings.Index(token, "/"); pos != -1 {
			if currentKey != "" {
				(*fp)[currentKey] = strings.TrimSpace(currentVal.String())
				currentVal.Reset()
			}
			currentKey = token[:pos]
			currentVal.WriteString(token[pos+1:] + " ")
		} else {
			currentVal.WriteString(token + " ")
		}
	}

	if currentKey != "" {
		(*fp)[currentKey] = strings.TrimSpace(currentVal.String())
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
