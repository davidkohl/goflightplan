package adexp

import "strings"

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
