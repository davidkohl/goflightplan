package goflightplan

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/davidkohl/goflightplan/adexp"
	"github.com/davidkohl/goflightplan/icao"
)

func Test_Parse_Fuzz(t *testing.T) {
	testSchema := []adexp.MessageSet{loadTestMessageSet(t)}
	pAdexp := adexp.NewParser(testSchema)
	pIcao := icao.NewParser(icao.ParserOpts{})

	dir := "./test/fpl/fuzz"
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Errorf("could not read directory: %v\n", err)
	}

	for _, file := range files {
		// Check if it's a regular file (not a directory)
		t.Run(file.Name(), func(t *testing.T) {
			// Read the test file
			if file.Type().IsRegular() {
				// Get the full path of the file
				filePath := dir + "/" + file.Name()
				fmt.Println("NOW PARSING:", file.Name())
				// Read the file's contents
				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("could not read file %s: %v", filePath, err)
				}
				//Print the contents of the file
				fpl, err := pAdexp.Parse(string(content))
				if err == nil {
					j, err := json.MarshalIndent(fpl, "", "\t")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("%v\n\n", string(j))
					fmt.Printf("%v\n", string(content))
					return
				}

				fpl, err = pIcao.Parse(string(content))
				if err == nil {
					j, err := json.MarshalIndent(fpl, "", "\t")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("----------------------\n\n%v\n\n", string(j))
					fmt.Printf("%v\n\n-----------------------------", string(content))
					return
				}
				t.Errorf("expected error to be nil after all parsers")

			}

		})

	}

}

func loadTestMessageSet(t *testing.T) adexp.MessageSet {
	t.Helper()
	messageSet, err := adexp.MessageSetFromJSON("./test/schema", "TestSet")
	if err != nil {
		t.Fatalf("Failed to load message set from JSON: %v", err)
	}
	return *messageSet
}
