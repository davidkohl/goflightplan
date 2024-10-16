package adexp

import (
	"fmt"
)

// parseStructuredField parses a structured field and returns its key, value (as a map), and any error
func (p *Parser) parseStructuredField(field DataField) (string, map[string]interface{}, error) {
	structuredData := make(map[string]interface{})

	for {
		subFieldName, subFieldValue, err := p.parseSubField(field.Subfields)
		if err != nil {
			return "", nil, err
		}
		if subFieldName == "" {
			break
		}
		structuredData[subFieldName] = subFieldValue
	}

	return field.DataItem, structuredData, nil
}

// parseSubField parses a subfield and returns its name, value, and any error
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

	subFieldDef := p.findField(subFieldName, subfields)
	if subFieldDef == nil {
		// This might be a new top-level field, so we need to backtrack
		p.currentPos -= len(subFieldName) + 2 // +2 for '-' and space
		return "", nil, nil
	}

	switch subFieldDef.Type {
	case Basicfield:
		_, value, err := p.parseBasicField(*subFieldDef)
		return subFieldName, value, err
	case StructuredField:
		_, value, err := p.parseStructuredField(*subFieldDef)
		return subFieldName, value, err
	default:
		return "", nil, fmt.Errorf("unsupported subfield type for '%s'", subFieldName)
	}
}
