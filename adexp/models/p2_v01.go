package models

import (
	"gitlab.com/davidkohl/goflightplan/adexp/schema"
)

type P2V01Model struct {
	TITLE   string
	REFDATA REFDATA
	ARCID   string
	SSRCODE string
	ADEP    string
	ADES    string
}

func (data *P2V01Model) Write(s string) error {
	for _, item := range schema.IcasV01.Items {
		var err error
		switch item.FRN {
		case 1:
			data.TITLE, err = parseBasicField(s, item.DataItem)
			if err != nil {
				return err
			}
		case 2:
			//data.REFDATA = parseREFDATA(s)
		case 3:
			data.ARCID, err = parseBasicField(s, item.DataItem)
			if err != nil {
				return err
			}
		case 4, 5, 6:
		}

	}
	return nil
}
