package icao

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Parse(t *testing.T) {
	// Load the test schema from JSON files
	parser := NewParser(ParserOpts{AFTNHeader: false})

	testCases := []struct {
		name        string
		filename    string
		description string
		expected    func(*testing.T, map[string]interface{})
	}{
		{
			name:     "BFD message short",
			filename: "CNL.txt",
			description: `
            This test verifies that fields that are not present in the schema are not parsed
            `,
			expected: func(t *testing.T, fpl map[string]interface{}) {
				if _, exists := fpl["WKTRC"]; exists {
					t.Errorf("Expected WKTRC to not be present, but it was")
				}
				if fpl["ARCID"] != "WMT912" {
					t.Errorf("Expected ARCID to be 'WMT912' but got %v\n", fpl["ARCID"])
				}
				if fpl["ADEP"] != "EDJA" {
					t.Errorf("Expected ADEP to be 'EDJA' but got %v\n", fpl["ARCID"])
				}
				if fpl["EOBT"] != "2010" {
					t.Errorf("Expected EOBT to be '2010' but got %v\n", fpl["ARCID"])
				}
				if fpl["ADES"] != "LIRF" {
					t.Errorf("Expected ADES to be 'LIRF' but got %v\n", fpl["LIRF"])
				}
				if fpl["DOF"] != "240228" {
					t.Errorf("Expected DOF to be '240228' but got %v\n", fpl["DOF"])
				}
			},
		},
		{
			name:        "FPL Basic",
			filename:    "FPL.txt",
			description: ``,
			expected: func(t *testing.T, fpl map[string]interface{}) {
				if fpl["WKTRC"] != "M" {
					t.Errorf("Expected WKTRC to be 'M' but got %v\n", fpl["WKTRC"])
				}
				if fpl["ADEP"] != "LPPR" {
					t.Errorf("Expected ADEP to be 'LPPR' but got %v\n", fpl["ADEP"])
				}
				if fpl["EOBT"] != "0600" {
					t.Errorf("Expected EOBT to be '0600' but got %v\n", fpl["EOBT"])
				}
				if fpl["ADES"] != "LFPG" {
					t.Errorf("Expected ADES to be 'LFPG' but got %v\n", fpl["LIRF"])
				}
				if fpl["EELT"] != "0155" {
					t.Errorf("Expected EELT to be '0155' but got %v\n", fpl["EELT"])
				}
				// Field 18
				if fpl["DOF"] != "060110" {
					t.Errorf("Expected DOF to be '060110' but got %v\n", fpl["DOF"])
				}
				if fpl["REG"] != "DESEL" {
					t.Errorf("Expected DOF to be 'DESEL' but got %v\n", fpl["DOF"])
				}
				if fpl["RMK"] != "THIS HAS SPACE AT END" {
					t.Errorf("Expected DOF to be 'DESEL' but got %v\n", fpl["DOF"])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read the test file
			content, err := os.ReadFile(filepath.Join("../test/fpl/icao", tc.filename))
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			// Parse the message
			result, err := parser.Parse(string(content))
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// Run the assertions
			tc.expected(t, result)
		})
	}
}

/*
fpl["TITLE"] = "FPL"
	fpl["ARCID"] = fields[0]             // Aircraft ID
	fpl["FLTRUL"] = string(fields[1][0]) // Flight rules
	fpl["FLTTYP"] = string(fields[1][1]) // Flight type
	fpl["ARCTYP"] = aircraftParts[0]     // Aircraft type
	fpl["WKTRC"] = aircraftParts[1]      // Wake turbulence category
	fpl["ADEP"] = departureInfo[:4]      // Departure aerodrome
	fpl["EOBT"] = departureInfo[4:]      // Estimated Off-Block Time (departure time)
	fpl["ROUTE"] = route                 // Route
	fpl["ADES"] = destinationInfo[:4]    // Destination aerodrome
	fpl["EELT"] = eelt                   // Estimated Elapsed Time

*/
