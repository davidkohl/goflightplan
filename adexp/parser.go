package adexp

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// Parser represents the ADEXP message parser
type Parser struct {
	MessageSet    []MessageSet
	buffer        bytes.Buffer
	currentPos    int
	message       string
	flightplan    map[string]interface{}
	currentSchema *StandardSchema
}

// NewParser creates a new Parser instance with the given schema
func NewParser(schema []MessageSet) *Parser {
	return &Parser{
		MessageSet: schema,
		flightplan: make(map[string]interface{}),
	}
}

// Parse parses the given ADEXP message and returns a map representation of the flight plan
func (p *Parser) Parse(message string) (map[string]interface{}, error) {
	p.currentPos = 0
	char, ok := validateMessage(message)
	if !ok {
		return nil, fmt.Errorf("invalid character '%v' found in message", string(char))
	}
	p.message = strings.ReplaceAll(message, "\n", " ")
	p.flightplan = make(map[string]interface{})

	titleStart := strings.Index(p.message, "-TITLE ")
	if titleStart == -1 {
		return nil, fmt.Errorf("TITLE field not found in the message")
	}
	titleStart += 7 // Length of "-TITLE "
	titleEnd := strings.Index(p.message[titleStart:], "-")
	if titleEnd == -1 {
		titleEnd = len(p.message)
	} else {
		titleEnd += titleStart
	}
	title := strings.TrimSpace(p.message[titleStart:titleEnd])

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

	return p.flightplan, nil
}

// parseNextField parses the next field in the message
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
	case ListField:
		return nil
		//return p.parseListField(field)
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

func validateMessage(message string) (rune, bool) {
	for _, char := range message {
		if c, ok := isValidCharacter(char); !ok {
			return c, false
		}
	}
	return rune(0), true
}

func isValidCharacter(char rune) (rune, bool) {
	// Upper case letters (A to Z)
	if unicode.IsUpper(char) {
		return rune(0), true
	}

	// Digits (0 to 9)
	if unicode.IsDigit(char) {
		return rune(0), true
	}

	// Special graphic characters
	switch char {
	case ' ', '(', ')', '-', '?', ':', '.', ',', '\'', '=', '+', '/':
		return rune(0), true
	}

	// Format effectors
	switch char {
	case '\r', '\n': // Carriage Return and Line Feed
		return rune(0), true
	}

	// Any other character is invalid
	return char, false
}
