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
	p.message = strings.ReplaceAll(message, "\n", " ")
	p.message = strings.TrimSpace(p.message)
	p.message = strings.TrimSuffix(p.message, "NNNN")
	p.flightplan = make(map[string]interface{})

	if err := p.validateMessage(); err != nil {
		return nil, err
	}

	if err := p.findTitle(); err != nil {
		return nil, err
	}

	for p.currentPos < len(p.message) {
		if err := p.parseNextField(); err != nil {
			return nil, fmt.Errorf("error parsing field: %w", err)
		}
	}

	return p.flightplan, nil
}

// validateMessage checks if the message contains only valid characters
func (p *Parser) validateMessage() error {
	for _, char := range p.message {
		if !isValidCharacter(char) {
			return fmt.Errorf("invalid character '%v' found in message", string(char))
		}
	}
	return nil
}

// findTitle locates the TITLE field and sets the appropriate schema
func (p *Parser) findTitle() error {
	titleStart := strings.Index(p.message, "-TITLE ")
	if titleStart == -1 {
		return fmt.Errorf("TITLE field not found in the message")
	}

	titleStart += 7 // Length of "-TITLE "
	titleEnd := strings.Index(p.message[titleStart:], "-")
	if titleEnd == -1 {
		titleEnd = len(p.message)
	} else {
		titleEnd += titleStart
	}
	title := strings.TrimSpace(p.message[titleStart:titleEnd])

	matchedSchema, err := p.findMatchingSchema(title)
	if err != nil {
		return err
	}
	p.currentSchema = matchedSchema
	return nil
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

	// Check if this is the start of a list field
	if fieldName == "BEGIN" {
		return p.handleListField()
	}

	field := p.findField(fieldName, p.currentSchema.Items)
	if field == nil {
		// If the field is not found in the schema, we'll skip it
		p.skipUnknownField()
		return nil
	}

	switch field.Type {
	case Basicfield:
		key, value, err := p.parseBasicField(*field)
		if err != nil {
			return err
		}
		p.flightplan[key] = value
	case StructuredField:
		key, value, err := p.parseStructuredField(*field)
		if err != nil {
			return err
		}
		p.flightplan[key] = value
	default:
		return fmt.Errorf("unknown field type for field '%s'", fieldName)
	}

	return nil
}

// findField finds a field in the given slice of DataFields
func (p *Parser) findField(fieldName string, fields []DataField) *DataField {
	for i := range fields {
		if fields[i].DataItem == fieldName {
			return &fields[i]
		}
	}
	return nil
}

// findMatchingSchema finds the matching schema for the given title
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

// skipUnknownField skips an unknown field in the message
func (p *Parser) skipUnknownField() {
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.currentPos++
	}
}

// isValidCharacter checks if a character is valid for ADEXP messages
func isValidCharacter(char rune) bool {
	if unicode.IsUpper(char) || unicode.IsDigit(char) {
		return true
	}

	switch char {
	case ' ', '(', ')', '-', '?', ':', '.', ',', '\'', '=', '+', '/', '\r', '\n':
		return true
	}

	return false
}

// handleListField handles the parsing of a list field
func (p *Parser) handleListField() error {
	// Skip the space after "BEGIN"
	p.currentPos++

	// Read the list field name
	p.buffer.Reset()
	for p.currentPos < len(p.message) && p.message[p.currentPos] != ' ' && p.message[p.currentPos] != '-' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}
	listFieldName := p.buffer.String()

	field := p.findField(listFieldName, p.currentSchema.Items)
	if field == nil || field.Type != ListField {
		return nil
	}

	key, value, err := p.parseListField(*field)
	if err != nil {
		return fmt.Errorf("error parsing list field '%s': %w", listFieldName, err)
	}
	p.flightplan[key] = value

	return nil
}
