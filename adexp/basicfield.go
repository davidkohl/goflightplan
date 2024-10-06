package adexp

import "strings"

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
