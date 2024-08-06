package schema

var IcasV01 = StandardUAP{
	Name:     "icas_0.1",
	Category: "BFD",
	Version:  0.1,
	Items: []DataField{
		{
			FRN:         1,
			DataItem:    "TITLE",
			Description: "Title of the ADEXP Message",
			Type:        Basicfield,
			Mendatory:   true,
		},
		{
			FRN:         2,
			DataItem:    "REFDATA",
			Description: "Message Reference with sender, receiver and sequence number",
			Type:        StructuredField,
			Mendatory:   true,
		},
		{
			FRN:         3,
			DataItem:    "ARCID",
			Description: "Aircraft id or callsign",
			Type:        Basicfield,
			Mendatory:   true,
		},
		{
			FRN:         4,
			DataItem:    "SSRCODE",
			Description: "Assigned SSRCODE",
			Type:        Basicfield,
		},
		{
			FRN:         5,
			DataItem:    "ADEP",
			Description: "Aerodrom of departure",
			Type:        Basicfield,
		},
		{
			FRN:         6,
			DataItem:    "ADES",
			Description: "Aerodrom of destination",
			Type:        Basicfield,
		},
		{
			FRN:         7,
			DataItem:    "ARCTYP",
			Description: "Aircraft type",
			Type:        Basicfield,
		},
		{
			FRN:         8,
			DataItem:    "IFPLID",
			Description: "Individual flight plan id",
			Type:        Basicfield,
		},
		{
			FRN:         9,
			DataItem:    "EOBT",
			Description: "Estimated off block time",
			Type:        Basicfield,
		},
		{
			FRN:         10,
			DataItem:    "ELDT",
			Description: "Estimated landing time",
			Type:        Basicfield,
		},
	},
}
