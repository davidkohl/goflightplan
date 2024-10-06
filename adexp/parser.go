package adexp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/davidkohl/goflightplan"
)

type Parser struct {
	MessageSet    []MessageSet
	buffer        bytes.Buffer
	currentPos    int
	message       string
	flightplan    *goflightplan.Flightplan
	fplwrapper    *goflightplan.FlightplanWrapper
	currentSchema *StandardSchema // New field to store the matched schema
}

func NewParser(schema []MessageSet) *Parser {
	return &Parser{
		MessageSet: schema,
		flightplan: &goflightplan.Flightplan{},
		fplwrapper: goflightplan.NewFlightplanWrapper(),
	}
}

func (p *Parser) Parse(message string) (*goflightplan.Flightplan, error) {

	p.currentPos = 0
	message = strings.ReplaceAll(message, "\n", " ")
	p.message = message

	titleStart := strings.Index(message, "-TITLE ")
	if titleStart == -1 {
		return nil, fmt.Errorf("TITLE field not found in the message")
	}
	titleStart += 7 // Length of "-TITLE "
	titleEnd := strings.Index(message[titleStart:], "-")
	if titleEnd == -1 {
		titleEnd = len(message)
	} else {
		titleEnd += titleStart
	}
	title := strings.TrimSpace(message[titleStart:titleEnd])

	// Find the matching schema
	matchedSchema, err := p.findMatchingSchema(title)
	if err != nil {
		return nil, err
	}
	p.currentSchema = matchedSchema

	for p.currentPos < len(p.message) {
		if err := p.parseNextField(); err != nil {
			if _, ok := err.(UnknownFieldError); ok {
				continue
			}
			return nil, fmt.Errorf("error parsing field: %w", err)
		}
	}
	fpl := &goflightplan.Flightplan{}
	fpl = p.flightplan
	p.flightplan = &goflightplan.Flightplan{}
	return fpl, nil
}

func (p *Parser) parseNextField() error {
	p.buffer.Reset()
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.currentPos++
	}

	if p.currentPos >= len(p.message) {
		return nil // End of message
	}

	p.currentPos++ // Skip the '-'
	for p.currentPos < len(p.message) && p.message[p.currentPos] != ' ' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}

	fieldName := p.buffer.String()
	field, _, err := p.findFieldInSchema(fieldName)
	if err != nil {
		// If the field is not found in the schema, we'll skip it
		p.skipUnknownField()
		return UnknownFieldError{FieldName: fieldName}
	}

	switch field.Type {
	case Basicfield:
		return p.parseBasicField(field)
	case StructuredField:
		return p.parseStructuredField(field)
	default:
		return fmt.Errorf("unknown field type for field '%s'", fieldName)
	}
}

func (p *Parser) skipUnknownField() {
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.currentPos++
	}
}

type UnknownFieldError struct {
	FieldName string
}

func (e UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field: %s", e.FieldName)
}

func (p *Parser) findFieldInSchema(fieldName string) (DataField, []DataField, error) {
	for _, item := range p.currentSchema.Items {
		if item.DataItem == fieldName {
			return item, item.Subfields, nil
		}
	}
	return DataField{}, nil, fmt.Errorf("field '%s' not found in schema", fieldName)
}

func (p *Parser) findMatchingSchema(title string) (*StandardSchema, error) {
	for _, messageSet := range p.MessageSet {
		for _, schema := range messageSet.Set {
			if schema.Category == title {
				return &schema, nil
			}
		}
	}
	return nil, fmt.Errorf("no matching schema found for title: %s", title)
}

func (p *Parser) setFieldValue(field DataField, value interface{}) error {
	v := reflect.ValueOf(p.flightplan).Elem()
	f := v.FieldByName(field.DataItem)

	if !f.IsValid() {
		return fmt.Errorf("field '%s' not found in Flightplan struct", field.DataItem)
	}

	if !f.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", field.DataItem)
	}

	switch f.Kind() {
	case reflect.String:
		f.SetString(value.(string))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// You might want to add proper parsing for integers
		return fmt.Errorf("integer parsing not implemented for field '%s'", field.DataItem)
	case reflect.Struct:
		// For structured fields, we need to marshal the map to JSON and then unmarshal into the struct
		mapValue, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map[string]interface{} for structured field '%s', got %T", field.DataItem, value)
		}
		jsonData, err := json.Marshal(mapValue)
		if err != nil {
			return fmt.Errorf("failed to marshal structured field '%s': %w", field.DataItem, err)
		}
		return json.Unmarshal(jsonData, f.Addr().Interface())
	default:
		fmt.Printf("Unhandled field type for %s: %v\n", field.DataItem, f.Kind())
	}

	return nil
}
