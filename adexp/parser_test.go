package adexp

import (
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/davidkohl/goflightplan"
)

func TestParse(t *testing.T) {
	// Load the default message set for testing
	messageSet, err := MessageSetFromJSON("../test/json", "SetName")
	if err != nil {
		t.Fatalf("Failed to load message schemas: %v", err)
	}

	parser := NewParser([]MessageSet{*messageSet}, ParserOpts{})

	// Define test cases
	testCases := []struct {
		name     string
		filename string
		expected func(*testing.T, *goflightplan.FlightplanWrapper)
	}{
		{
			name:     "BFD message",
			filename: "007_BFD.txt",
			expected: func(t *testing.T, fp *goflightplan.FlightplanWrapper) {
				if fp.Flightplan.TITLE != "BFD" {
					t.Errorf("Expected TITLE to be BFD, got %s", fp.Flightplan.TITLE)
				}
				if fp.Flightplan.ARCID != "WMT3GH" {
					t.Errorf("Expected ARCID to be WMT3GH, got %s", fp.Flightplan.ARCID)
				}
				// Add more checks based on the expected content
			},
		},
		{
			name:     "CFD message",
			filename: "008_CFD.txt",
			expected: func(t *testing.T, fp *goflightplan.FlightplanWrapper) {
				if fp.Flightplan.TITLE != "CFD" {
					t.Errorf("Expected TITLE to be ACT, got %s", fp.Flightplan.TITLE)
				}
				if fp.Flightplan.ARCID != "WMT3GH" {
					t.Errorf("Expected ARCID to be 'WMT3GH', got '%s'", fp.Flightplan.ARCID)
				}
				// Add more checks based on the expected content
			},
		},
		{
			name:     "TFD message",
			filename: "009_TFD.txt",
			expected: func(t *testing.T, fp *goflightplan.FlightplanWrapper) {
				if fp.Flightplan.TITLE != "TFD" {
					t.Errorf("Expected TITLE to be ACT, got %s", fp.Flightplan.TITLE)
				}
				if fp.Flightplan.ARCID != "WZZ70BK" {
					t.Errorf("Expected ARCID to be 'WMT3GH', got '%s'", fp.Flightplan.ARCID)
				}
				// Add more checks based on the expected content
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read the test file
			content, err := os.ReadFile(filepath.Join("../test/fpl/adexp", tc.filename))
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
