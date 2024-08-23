package icasv01

import (
	"gitlab.com/davidkohl/goflightplan/adexp"
)

func init() {
	MessageSet["BFD"] = BFD
}

var BFD = adexp.StandardSchema{
	Name:     "icas_0.1",
	Category: "BFD",
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
		{
			FRN:         11,
			DataItem:    "WKTRC",
			Description: "Wake Turbulance Category",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         12,
			DataItem:    "EOBD",
			Description: "Day of Flight",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         13,
			DataItem:    "FLTTYP",
			Description: "Day of Flight",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         14,
			DataItem:    "FLTRUL",
			Description: "Day of Flight",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "FPLCAT",
			Description: "FPLCAT",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "CTOT",
			Description: "CTOT",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "SID",
			Description: "SID",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "STAR",
			Description: "STAR",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "REG",
			Description: "REG",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "NEWCTOT",
			Description: "NEWCTOT",
			Type:        adexp.Basicfield,
		},
		{
			FRN:         15,
			DataItem:    "ROUTE",
			Description: "ROUTE",
			Type:        adexp.Basicfield,
		},
	},
}
