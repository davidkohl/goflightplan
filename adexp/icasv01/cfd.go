package icasv01

import "gitlab.com/davidkohl/goflightplan/adexp"

func init() {
	MessageSet["CFD"] = CFD
}

var CFD = adexp.StandardSchema{
	Name:     "icas_0.1",
	Category: "CFD",
	Version:  0.1,
	Items: []adexp.DataField{
		{
			FRN:         1,
			DataItem:    "TITLE",
			Description: "Title of the ADEXP Message",
			Type:        adexp.Basicfield,
			Mendatory:   true,
		},
		{
			FRN:         2,
			DataItem:    "REFDATA",
			Description: "Message Reference with sender, receiver and sequence number",
			Type:        adexp.StructuredField,
			Mendatory:   true,
		},
		{
			FRN:         3,
			DataItem:    "ARCID",
			Description: "Aircraft id or callsign",
			Type:        adexp.Basicfield,
			Mendatory:   true,
		},
		{
			FRN:         4,
			DataItem:    "SSRCODE",
			Description: "Assigned SSRCODE",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         5,
			DataItem:    "ADEP",
			Description: "Aerodrom of departure",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         6,
			DataItem:    "ADES",
			Description: "Aerodrom of destination",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         7,
			DataItem:    "ARCTYP",
			Description: "Aircraft type",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         8,
			DataItem:    "IFPLID",
			Description: "Individual flight plan id",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         9,
			DataItem:    "EOBT",
			Description: "Estimated off block time",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         10,
			DataItem:    "ELDT",
			Description: "Estimated landing time",
			Type:        adexp.Basicfield,
		},
	},
}
