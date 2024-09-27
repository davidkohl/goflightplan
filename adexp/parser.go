package adexp

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"gitlab.com/davidkohl/goflightplan"
)

type Parser struct {
	schema     []MessageSet
	buffer     bytes.Buffer
	currentPos int
	message    string
	flightplan *goflightplan.Flightplan
	fplwrapper *goflightplan.FlightplanWrapper
}

func NewParser(schema []MessageSet) *Parser {
	return &Parser{
		schema:     schema,
		flightplan: &goflightplan.Flightplan{},
		fplwrapper: goflightplan.NewFlightplanWrapper(),
	}
}

func (p *Parser) Parse(message string) (*goflightplan.FlightplanWrapper, error) {
	p.message = message
	p.currentPos = 0

	for p.currentPos < len(p.message) {
		if err := p.parseNextField(); err != nil {
			// If it's an unknown field error, we'll just continue
			if _, ok := err.(UnknownFieldError); ok {
				continue
			}
			// For other types of errors, we'll return the error
			return nil, fmt.Errorf("error parsing field: %w", err)
		}
	}

	p.fplwrapper.Flightplan = *p.flightplan
	p.fplwrapper.Raw = message
	return p.fplwrapper, nil
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
	field, err := p.findFieldInSchema(fieldName)
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

// Add this new method to skip unknown fields
func (p *Parser) skipUnknownField() {
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.currentPos++
	}
}

// Define a custom error type for unknown fields
type UnknownFieldError struct {
	FieldName string
}

func (e UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field: %s", e.FieldName)
}
func (p *Parser) parseBasicField(field DataField) error {
	p.buffer.Reset()
	p.currentPos++ // Skip the space after field name
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}

	value := strings.TrimSpace(p.buffer.String())
	return p.setFieldValue(field, value)
}

func (p *Parser) parseStructuredField(field DataField) error {
	// For simplicity, we'll treat structured fields as basic fields for now
	// In a real implementation, you'd recursively parse subfields here
	return p.parseBasicField(field)
}

func (p *Parser) findFieldInSchema(fieldName string) (DataField, error) {
	for _, messageSet := range p.schema {
		for _, schema := range messageSet.Set {
			for _, item := range schema.Items {
				if item.DataItem == fieldName {
					return item, nil
				}
			}
		}
	}
	return DataField{}, fmt.Errorf("field '%s' not found in schema", fieldName)
}

func (p *Parser) setFieldValue(field DataField, value string) error {
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
		f.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// You might want to add proper parsing for integers
		return fmt.Errorf("integer parsing not implemented for field '%s'", field.DataItem)
	// Add more types as needed
	default:
		return nil //fmt.Errorf("unsupported type for field '%s'", field.DataItem)
	}

	return nil
}
