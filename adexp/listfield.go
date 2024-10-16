package adexp

import (
	"fmt"
	"log"
)

// parseListField parses a list field and returns its key, value (as a slice), and any error
func (p *Parser) parseListField(field DataField) (string, []interface{}, error) {
	listData := make([]interface{}, 0)

	for {
		item, err := p.parseListItem(field.Subfields)
		if err != nil {
			return "", nil, fmt.Errorf("error parsing list item: %w", err)
		}
		listData = append(listData, item)

		// Check for END marker after each item
		if p.checkForEndMarker(field.DataItem) {
			break
		}
	}

	return field.DataItem, listData, nil
}

// parseListItem parses a single item in a list field
func (p *Parser) parseListItem(subfields []DataField) (interface{}, error) {
	if len(subfields) == 1 && subfields[0].Type == Basicfield {
		// Handle simple list (e.g., EQCST)
		return p.parseSimpleListItem(subfields[0])
	}

	// Handle structured list (e.g., RTEPTS)
	return p.parseStructuredListItem(subfields)
}

// parseSimpleListItem parses a single item in a simple list field
func (p *Parser) parseSimpleListItem(subfield DataField) (string, error) {
	subFieldName, subFieldValue, err := p.parseSubField([]DataField{subfield})
	if err != nil {
		return "", fmt.Errorf("error parsing simple list item: %w", err)
	}
	if subFieldName == "" {
		return "", nil // End of list
	}
	return subFieldValue.(string), nil
}

// parseStructuredListItem parses a single item in a structured list field
func (p *Parser) parseStructuredListItem(subfields []DataField) (map[string]interface{}, error) {
	item := make(map[string]interface{})

	for {
		subFieldName, subFieldValue, err := p.parseSubField(subfields)
		if err != nil {
			return nil, fmt.Errorf("error parsing structured list item: %w", err)
		}
		if subFieldName == "" {
			break // End of current item or start of next item
		}

		// If the subfield value is a map, flatten it
		if subValue, ok := subFieldValue.(map[string]interface{}); ok {
			for k, v := range subValue {
				item[k] = v
			}
		} else {
			item[subFieldName] = subFieldValue
		}

		// Check if we've reached the start of a new item
		if p.isStartOfNewItem(subfields) {
			break
		}
	}

	if len(item) == 0 {
		return nil, nil // No item parsed, end of list
	}

	return item, nil
}

// isStartOfNewItem checks if the current position is the start of a new item in the list
func (p *Parser) isStartOfNewItem(subfields []DataField) bool {
	originalPos := p.currentPos
	p.buffer.Reset()

	// Skip spaces
	for p.currentPos < len(p.message) && p.message[p.currentPos] == ' ' {
		p.currentPos++
	}

	// Check for '-' character
	if p.currentPos < len(p.message) && p.message[p.currentPos] == '-' {
		p.currentPos++
		// Read the field name
		for p.currentPos < len(p.message) && p.message[p.currentPos] != ' ' {
			p.buffer.WriteByte(p.message[p.currentPos])
			p.currentPos++
		}
		fieldName := p.buffer.String()

		// Check if this field name is in the subfields list
		for _, subfield := range subfields {
			if subfield.DataItem == fieldName {
				p.currentPos = originalPos
				return true
			}
		}
	}

	p.currentPos = originalPos
	return false
}

// checkForEndMarker checks if the next field is the END marker for the list
func (p *Parser) checkForEndMarker(fieldName string) bool {
	originalPos := p.currentPos
	p.buffer.Reset()

	// Skip spaces
	for p.currentPos < len(p.message) && p.message[p.currentPos] == ' ' {
		p.currentPos++
	}

	if p.currentPos >= len(p.message) || p.message[p.currentPos] != '-' {
		p.currentPos = originalPos
		return false
	}
	p.currentPos++ // Skip the '-'

	for p.currentPos < len(p.message) && p.message[p.currentPos] != ' ' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}

	if p.buffer.String() != "END" {
		p.currentPos = originalPos
		return false
	}

	// Consume the field name after END
	p.buffer.Reset()
	p.currentPos++ // Skip the space after END
	for p.currentPos < len(p.message) && p.message[p.currentPos] != ' ' && p.message[p.currentPos] != '-' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}

	endFieldName := p.buffer.String()
	if endFieldName != fieldName {
		log.Printf("Warning: END field name '%s' does not match BEGIN field name '%s'", endFieldName, fieldName)
	}

	return true
}
