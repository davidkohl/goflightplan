package adexp

import (
	"strings"
)

// parseBasicField parses a basic field and returns its key, value, and any error
func (p *Parser) parseBasicField(field DataField) (string, string, error) {
	p.buffer.Reset()
	p.currentPos++ // Skip the space after field name
	for p.currentPos < len(p.message) && p.message[p.currentPos] != '-' {
		p.buffer.WriteByte(p.message[p.currentPos])
		p.currentPos++
	}

	value := strings.TrimSpace(p.buffer.String())
	return field.DataItem, value, nil
}
