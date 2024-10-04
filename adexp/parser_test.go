package adexp

import (
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/davidkohl/goflightplan"
)

func TestParser_Parse(t *testing.T) {
	// Load the test schema from JSON files
	testSchema := []MessageSet{loadTestMessageSet(t)}
	parser := NewParser(testSchema)

	testCases := []struct {
		name     string
		filename string
		expected func(*testing.T, *goflightplan.Flightplan)
	}{
		{
			name:     "BFD message",
			filename: "007_BFD.txt",
			expected: func(t *testing.T, fp *goflightplan.Flightplan) {
				if fp.TITLE != "BFD" {
					t.Errorf("Expected TITLE to be BFD, got %s", fp.TITLE)
				}
				if fp.ARCID != "DLH151" {
					t.Errorf("Expected ARCID to be DLH151, got %s", fp.ARCID)
				}
				// Add more checks based on the expected content of BFD message
				if fp.ADEP != "EDDW" {
					t.Errorf("Expected ADEP to be EDDW, got %s", fp.ADEP)
				}
				if fp.ADES != "GMME" {
					t.Errorf("Expected ADES to be GMME, got %s", fp.ADES)
				}
				if fp.REFDATA.SENDER.FAC != "EBBUZXZQ" {
					t.Errorf("Expected SENDER to be EBBUZXZQ, got %s", fp.REFDATA.SENDER.FAC)
				}
				if fp.REFDATA.RECVR.FAC != "EBSZZXZQ" {
					t.Errorf("Expected RECVR to be EBSZZXZQ, got %s", fp.REFDATA.RECVR.FAC)
				}
				if fp.REFDATA.SEQNUM != "006" {
					t.Errorf("Expected SEQNUM to be 006, got %s", fp.REFDATA.SEQNUM)
				}

				// Check other fields specific to BFD message
			},
		},
		{
			name:     "CFD message",
			filename: "008_CFD.txt",
			expected: func(t *testing.T, fp *goflightplan.Flightplan) {
				if fp.TITLE != "CFD" {
					t.Errorf("Expected TITLE to be CFD, got %s", fp.TITLE)
				}
				if fp.ARCID != "DLH151" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%s'", fp.ARCID)
				}
				// Add more checks based on the expected content of CFD message
			},
		},
		{
			name:     "TFD message",
			filename: "009_TFD.txt",
			expected: func(t *testing.T, fp *goflightplan.Flightplan) {
				if fp.TITLE != "TFD" {
					t.Errorf("Expected TITLE to be TFD, got %s", fp.TITLE)
				}
				if fp.ARCID != "DLH151" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%s'", fp.ARCID)
				}
				// Add more checks based on the expected content of TFD message
			},
		},
		// Add more test cases for other message types
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

// Helper function to load the MessageSet
func loadTestMessageSet(t *testing.T) MessageSet {
	t.Helper()
	messageSet, err := MessageSetFromJSON("../test/schema", "TestSet")
	if err != nil {
		t.Fatalf("Failed to load message set from JSON: %v", err)
	}
	return *messageSet
}
