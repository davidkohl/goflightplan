package adexp

// parseListField parses a list field and adds it to the flightplan map
func (p *Parser) parseListField(field DataField) error {
	listData := make([]map[string]interface{}, 0)
	for p.currentPos < len(p.message) {
		itemData := make(map[string]interface{})
		for _, subField := range field.Subfields {
			subFieldName, subFieldValue, err := p.parseSubField([]DataField{subField})
			if err != nil {
				return err
			}
			if subFieldName == "" {
				break
			}
			itemData[subFieldName] = subFieldValue
		}
		if len(itemData) > 0 {
			listData = append(listData, itemData)
		} else {
			break
		}
	}

	p.flightplan[field.DataItem] = listData
	return nil
}
