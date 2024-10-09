package adexp

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_Parse(t *testing.T) {
	// Load the test schema from JSON files
	testSchema := []MessageSet{loadTestMessageSet(t)}
	parser := NewParser(testSchema)

	testCases := []struct {
		name        string
		filename    string
		description string
		expected    func(*testing.T, map[string]interface{})
	}{
		{
			name:     "BFD message short",
			filename: "BFD_short.txt",
			description: `
            This test verifies that fields that are not present in the schema are not parsed
            `,
			expected: func(t *testing.T, fp map[string]interface{}) {
				if _, exists := fp["WKTRC"]; exists {
					t.Errorf("Expected WKTRC to not be present, but it was")
				}
			},
		},
		{
			name:     "BFD message",
			filename: "BFD.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "BFD" {
					t.Errorf("Expected TITLE to be BFD, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "DLH151" {
					t.Errorf("Expected ARCID to be DLH151, got %v", fp["ARCID"])
				}
				if fp["ADEP"] != "EDDW" {
					t.Errorf("Expected ADEP to be EDDW, got %v", fp["ADEP"])
				}
				if fp["ADES"] != "GMME" {
					t.Errorf("Expected ADES to be GMME, got %v", fp["ADES"])
				}
				refdata, ok := fp["REFDATA"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected REFDATA to be a map[string]interface{}")
				} else {
					sender, ok := refdata["SENDER"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected SENDER to be a map[string]interface{}")
					} else if sender["FAC"] != "EBBUZXZQ" {
						t.Errorf("Expected SENDER.FAC to be EBBUZXZQ, got %v", sender["FAC"])
					}
					recvr, ok := refdata["RECVR"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected RECVR to be a map[string]interface{}")
					} else if recvr["FAC"] != "EBSZZXZQ" {
						t.Errorf("Expected RECVR.FAC to be EBSZZXZQ, got %v", recvr["FAC"])
					}
					if refdata["SEQNUM"] != "006" {
						t.Errorf("Expected SEQNUM to be 006, got %v", refdata["SEQNUM"])
					}
				}
				if fp["WKTRC"] != "M" {
					t.Errorf("Expected WKTRC to be 'M', got %v", fp["WKTRC"])
				}
			},
		},
		{
			name:     "CFD message",
			filename: "CFD.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "CFD" {
					t.Errorf("Expected TITLE to be CFD, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "DLH151" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%v'", fp["ARCID"])
				}
			},
		},
		{
			name:     "TFD message",
			filename: "TFD.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "TFD" {
					t.Errorf("Expected TITLE to be TFD, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "DLH151" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%v'", fp["ARCID"])
				}
			},
		},
		{
			name:     "SAM message",
			filename: "SAM.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "SAM" {
					t.Errorf("Expected TITLE to be SAM, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "AMC101" {
					t.Errorf("Expected ARCID to be 'AMC101', got '%v'", fp["ARCID"])
				}
				if fp["ADEP"] != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%v'", fp["ADEP"])
				}
				if fp["ADES"] != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%v'", fp["ADES"])
				}
				if fp["EOBD"] != "160224" {
					t.Errorf("Expected EOBD to be '160224', got '%v'", fp["EOBD"])
				}
				if fp["EOBT"] != "0945" {
					t.Errorf("Expected EOBT to be '0945', got '%v'", fp["EOBT"])
				}
				if fp["CTOT"] != "1200" {
					t.Errorf("Expected CTOT to be '1200', got '%v'", fp["CTOT"])
				}
				if fp["REGUL"] != "LMMLA24" {
					t.Errorf("Expected REGUL to be 'LMMLA24', got '%v'", fp["REGUL"])
				}
				tto, ok := fp["TTO"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected TTO to be a map[string]interface{}")
				} else {
					if tto["PTID"] != "GZO" {
						t.Errorf("Expected TTO.PTID to be 'GZO', got '%v'", tto["PTID"])
					}
					if tto["TO"] != "1438" {
						t.Errorf("Expected TTO.TO to be '1438', got '%v'", tto["TO"])
					}
					if tto["FL"] != "F060" {
						t.Errorf("Expected TTO.FL to be 'F060', got '%v'", tto["FL"])
					}
				}
				if fp["TAXITIME"] != "0010" {
					t.Errorf("Expected TAXITIME to be '0010', got '%v'", fp["TAXITIME"])
				}
				if fp["REGCAUSE"] != "WA 84" {
					t.Errorf("Expected REGCAUSE to be 'WA 84', got '%v'", fp["REGCAUSE"])
				}
				if fp["RVR"] != "100" {
					t.Errorf("Expected RVR to be '100', got '%v'", fp["RVR"])
				}
			},
		},
		{
			name:     "SRM message",
			filename: "SRM.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "SRM" {
					t.Errorf("Expected TITLE to be SRM, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "AMC101" {
					t.Errorf("Expected ARCID to be 'AMC101', got '%v'", fp["ARCID"])
				}
				if fp["ADEP"] != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%v'", fp["ADEP"])
				}
				if fp["ADES"] != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%v'", fp["ADES"])
				}
				if fp["EOBD"] != "160224" {
					t.Errorf("Expected EOBD to be '160224', got '%v'", fp["EOBD"])
				}
				if fp["EOBT"] != "0945" {
					t.Errorf("Expected EOBT to be '0945', got '%v'", fp["EOBT"])
				}
				if fp["NEWCTOT"] != "1200" {
					t.Errorf("Expected NEWCTOT to be '1200', got '%v'", fp["NEWCTOT"])
				}
				if fp["REGUL"] != "LMMLA24" {
					t.Errorf("Expected REGUL to be 'LMMLA24', got '%v'", fp["REGUL"])
				}
				tto, ok := fp["TTO"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected TTO to be a map[string]interface{}")
				} else {
					if tto["PTID"] != "GZO" {
						t.Errorf("Expected TTO.PTID to be 'GZO', got '%v'", tto["PTID"])
					}
					if tto["TO"] != "1438" {
						t.Errorf("Expected TTO.TO to be '1438', got '%v'", tto["TO"])
					}
					if tto["FL"] != "F060" {
						t.Errorf("Expected TTO.FL to be 'F060', got '%v'", tto["FL"])
					}
				}
				if fp["TAXITIME"] != "0010" {
					t.Errorf("Expected TAXITIME to be '0010', got '%v'", fp["TAXITIME"])
				}
				if fp["REGCAUSE"] != "WA 84" {
					t.Errorf("Expected REGCAUSE to be 'WA 84', got '%v'", fp["REGCAUSE"])
				}
			},
		},
		{
			name:     "SLC message",
			filename: "SLC.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "SLC" {
					t.Errorf("Expected TITLE to be SLC, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "AMC101" {
					t.Errorf("Expected ARCID to be 'AMC101', got '%v'", fp["ARCID"])
				}
				if fp["ADEP"] != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%v'", fp["ADEP"])
				}
				if fp["ADES"] != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%v'", fp["ADES"])
				}
				if fp["EOBD"] != "080901" {
					t.Errorf("Expected EOBD to be '080901', got '%v'", fp["EOBD"])
				}
				if fp["EOBT"] != "0945" {
					t.Errorf("Expected EOBT to be '0945', got '%v'", fp["EOBT"])
				}
				if fp["REASON"] != "VOID" {
					t.Errorf("Expected REASON to be 'VOID', got '%v'", fp["REASON"])
				}
				if fp["TAXITIME"] != "0020" {
					t.Errorf("Expected TAXITIME to be '0020', got '%v'", fp["TAXITIME"])
				}
				if fp["COMMENT"] != "FLIGHT CANCELLED" {
					t.Errorf("Expected COMMENT to be 'FLIGHT CANCELLED', got '%v'", fp["COMMENT"])
				}
			},
		},
		{
			name:     "FLS message",
			filename: "FLS.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "FLS" {
					t.Errorf("Expected TITLE to be FLS, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "AMC101" {
					t.Errorf("Expected ARCID to be 'AMC101', got '%v'", fp["ARCID"])
				}
				if fp["ADEP"] != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%v'", fp["ADEP"])
				}
				if fp["ADES"] != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%v'", fp["ADES"])
				}
				if fp["EOBD"] != "080901" {
					t.Errorf("Expected EOBD to be '080901', got '%v'", fp["EOBD"])
				}
				if fp["EOBT"] != "0945" {
					t.Errorf("Expected EOBT to be '0945', got '%v'", fp["EOBT"])
				}
				if fp["TAXITIME"] != "0020" {
					t.Errorf("Expected TAXITIME to be '0020', got '%v'", fp["TAXITIME"])
				}
				if fp["COMMENT"] != "RVR UNKNOWN" {
					t.Errorf("Expected COMMENT to be 'RVR UNKNOWN', got '%v'", fp["COMMENT"])
				}
				if fp["REGCAUSE"] != "WA 84" {
					t.Errorf("Expected REGCAUSE to be 'WA 84', got '%v'", fp["REGCAUSE"])
				}
				if fp["REGUL"] != "UZZU11" {
					t.Errorf("Expected REGUL to be 'UZZU11', got '%v'", fp["REGUL"])
				}
			},
		},
		{
			name:     "DES message",
			filename: "DES.txt",
			expected: func(t *testing.T, fp map[string]interface{}) {
				if fp["TITLE"] != "DES" {
					t.Errorf("Expected TITLE to be DES, got %v", fp["TITLE"])
				}
				if fp["ARCID"] != "AMC101" {
					t.Errorf("Expected ARCID to be 'AMC101', got '%v'", fp["ARCID"])
				}
				if fp["ADEP"] != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%v'", fp["ADEP"])
				}
				if fp["ADES"] != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%v'", fp["ADES"])
				}
				if fp["EOBD"] != "080901" {
					t.Errorf("Expected EOBD to be '080901', got '%v'", fp["EOBD"])
				}
				if fp["EOBT"] != "0945" {
					t.Errorf("Expected EOBT to be '0945', got '%v'", fp["EOBT"])
				}
				if fp["TAXITIME"] != "0020" {
					t.Errorf("Expected TAXITIME to be '0020', got '%v'", fp["TAXITIME"])
				}
				if fp["COMMENT"] != "NEW ATFM MESSAGES MAY POSSIBLY BE PUBLISHED AT 2 HOURS BEFORE THE EOBT" {
					t.Errorf("Expected COMMENT to be 'NEW ATFM MESSAGES MAY POSSIBLY BE PUBLISHED AT 2 HOURS BEFORE THE EOBT', got '%v'", fp["COMMENT"])
				}
			},
		},
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

func Test_Parse_Errors(t *testing.T) {
	// Load the test schema from JSON files
	testSchema := []MessageSet{loadTestMessageSet(t)}
	parser := NewParser(testSchema)

	testCases := []struct {
		name        string
		filename    string
		description string
		expectError bool
		expected    func(*testing.T, error)
	}{
		{
			name:        "ADEXP NO TITLE",
			filename:    "ADEXP_no_title.txt",
			expectError: true,
			expected: func(t *testing.T, err error) {
				if err == nil || err.Error() != "TITLE field not found in the message" {
					t.Errorf("Expected error 'TITLE field not found in the message', got: %v", err)
				}
			},
		},
		{
			name:        "NO SCHEMA",
			filename:    "ADEXP_no_schema.txt",
			expectError: true,
			expected: func(t *testing.T, err error) {
				if err == nil || err.Error() != "no matching schema found for title: ABC" {
					t.Errorf("Expected error 'no matching schema found for title: ABC', got: %v", err)
				}
			},
		},
		{
			name:        "EMPTY MSG",
			filename:    "ADEXP_empty.txt",
			expectError: true,
			expected: func(t *testing.T, err error) {
				if err == nil || err.Error() != "TITLE field not found in the message" {
					t.Errorf("Expected error 'TITLE field not found in the message', got: %v", err)
				}
			},
		},
		{
			name:        "INVALID CHAR",
			filename:    "ADEXP_invalid_char.txt",
			expectError: true,
			expected: func(t *testing.T, err error) {
				// The behavior for invalid characters might need to be defined
				if err == nil {
					t.Errorf("Expected an error for invalid characters, got nil")
				}
				fmt.Println(err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read the test file
			content, err := os.ReadFile(filepath.Join("../test/fpl/adexp", tc.filename))
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			// Parse the message
			_, err = parser.Parse(string(content))
			if err == nil && tc.expectError {
				t.Fatalf("Expected an error, but got nil")
			}
			if err != nil && !tc.expectError {
				t.Fatalf("Did not expect an error, but got: %v", err)
			}

			// Run the assertions
			if tc.expectError {
				tc.expected(t, err)
			}
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
