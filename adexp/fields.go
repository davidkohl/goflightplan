package adexp

import "fmt"

const (
	Basicfield = iota
	ListField
	StructuredField
)

type StandardSchema struct {
	Name     string
	Category string
	Version  float64
	Items    []DataField
}

// DataField describes FRN(Field Reference Number)
type DataField struct {
	FRN         uint8
	DataItem    string
	Description string
	Type        uint8
	Mendatory   bool
}

type ADEXPModel interface {
	Write(s string) error
}

type MessageSet struct {
	Name string
	Set  map[string]StandardSchema
}

var ErrorFieldNotPresent = fmt.Errorf("field not present")
