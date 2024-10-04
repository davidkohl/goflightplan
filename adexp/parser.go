package adexp

import (
	"bytes"
	"encoding/json"
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

func (p *Parser) Parse(message string) (*goflightplan.Flightplan, error) {

	p.currentPos = 0
	message = strings.ReplaceAll(message, "\n", " ")
	p.message = message
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

	structuredData := make(map[string]interface{})
	for p.currentPos < len(p.message) {
		subFieldName, subFieldValue, err := p.parseSubField(field.Subfields)
		if err != nil {
			return err
		}
		if subFieldName == "" {
			break
		}
		structuredData[subFieldName] = subFieldValue
	}

	return p.setFieldValue(field, structuredData)
}

func (p *Parser) parseSubField(subfields []DataField) (string, interface{}, error) {
	p.buffer.Reset()
	for p.currentPos < len(p.message) && p.message[p.currentPos] == ' ' {
		p.currentPos++
	}
	if p.currentPos >= len(p.message) || p.message[p.currentPos] != '-' {
		return "", nil, nil // End of structured field
	}
	p.currentPos++ // Skip the '-'

	for p.currentPos < len(p.message) && p.message[p.currentPos] != ' ' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}
	subFieldName := p.buffer.String()

	p.currentPos++ // Skip the space after the field name
	p.buffer.Reset()

	// Find the subfield definition
	var subFieldDef *DataField
	for i := range subfields {
		if subfields[i].DataItem == subFieldName {
			subFieldDef = &subfields[i]
			break
		}
	}

	if subFieldDef == nil {
		// This might be a new top-level field, so we need to backtrack
		p.currentPos -= len(subFieldName) + 2 // +2 for '-' and space
		return "", nil, nil
	}

	// Handle nested structures
	if subFieldDef.Type == StructuredField {
		nestedData, err := p.parseNestedStructure(subFieldDef.Subfields)
		if err != nil {
			return "", nil, err
		}
		return subFieldName, nestedData, nil
	}

	// Parse simple value
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}
	subFieldValue := strings.TrimSpace(p.buffer.String())

	return subFieldName, subFieldValue, nil
}

func (p *Parser) parseNestedStructure(subfields []DataField) (map[string]interface{}, error) {
	nestedData := make(map[string]interface{})
	for {
		subFieldName, subFieldValue, err := p.parseSubField(subfields)
		if err != nil {
			return nil, err
		}
		if subFieldName == "" {
			break
		}
		nestedData[subFieldName] = subFieldValue
	}
	return nestedData, nil
}

func (p *Parser) findFieldInSchema(fieldName string) (DataField, []DataField, error) {
	for _, messageSet := range p.schema {
		for _, schema := range messageSet.Set {
			for _, item := range schema.Items {
				if item.DataItem == fieldName {
					return item, item.Subfields, nil
				}
			}
		}
	}
	return DataField{}, nil, fmt.Errorf("field '%s' not found in schema", fieldName)
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
