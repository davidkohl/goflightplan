package icao

import "fmt"

type ICAOModel interface {
	Write(s string) error
}

var ErrorFieldNotPresent = fmt.Errorf("field not present")
