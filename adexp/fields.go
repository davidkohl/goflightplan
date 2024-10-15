package adexp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	Basicfield = iota
	ListField
	StructuredField
)

type StandardSchema struct {
	Name     string
	Category string
	Version  string
	Items    []DataField
}

// DataField describes FRN(Field Reference Number)
type DataField struct {
	FRN         uint8
	DataItem    string
	Description string
	Type        uint8
	Mendatory   bool
	Target      string
	Subfields   []DataField
}

type MessageSet struct {
	Name string
	Set  map[string]StandardSchema
}

var ErrorFieldNotPresent = fmt.Errorf("field not present")
var ErrorMendatory = fmt.Errorf("mendatory field not present")

func MessageSetFromJSON(p string, n string) (*MessageSet, error) {
	var set MessageSet
	set.Set = make(map[string]StandardSchema, 0)
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		var schema = StandardSchema{}
		// Check if it's a regular file (not a directory)
		if file.Type().IsRegular() {
			// Get the full path of the file
			if strings.Split(file.Name(), ".")[1] != "json" {
				continue
			}
			filePath := p + "/" + file.Name()

			// Read the file's contents
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(content, &schema)
			if err != nil {
				return nil, err
			}
			set.Set[schema.Category] = schema
		}
	}

	if len(set.Set) == 0 {
		return nil, errors.New("length of set is 0")
	}

	return &set, nil
}
