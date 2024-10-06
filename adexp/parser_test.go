package adexp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davidkohl/goflightplan"
)

func Test_Parse(t *testing.T) {
	// Load the test schema from JSON files
	testSchema := []MessageSet{loadTestMessageSet(t)}
	parser := NewParser(testSchema)

	testCases := []struct {
		name        string
		filename    string
		description string
		expected    func(*testing.T, *goflightplan.Flightplan)
	}{{
		name:     "BFD message short",
		filename: "007_BFD_short.txt",
		description: `
		This test verifies that fields that are not present in the schema are not parsed
		`,
		expected: func(t *testing.T, fp *goflightplan.Flightplan) {
			if fp.WKTRC != "" {
				t.Errorf("Expected WKTRC to be empty, got value %s", fp.WKTRC)
			}

			// Check other fields specific to BFD message
		},
	},
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
				if fp.WKTRC != "M" {
					t.Errorf("Expected WKTRC to be 'M', got %s", fp.WKTRC)
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
		{
			/*
				Test all fields according to the ATFCM USERS MANUAL v.22.1
			*/
			name:     "SAM message",
			filename: "SAM.txt",
			expected: func(t *testing.T, fp *goflightplan.Flightplan) {
				if fp.TITLE != "SAM" {
					t.Errorf("Expected TITLE to be TFD, got %s", fp.TITLE)
				}
				if fp.ARCID != "AMC101" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%s'", fp.ARCID)
				}
				if fp.ADEP != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%s'", fp.ADEP)
				}
				if fp.ADES != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%s'", fp.ADES)
				}
				if fp.EOBD != "160224" {
					t.Errorf("Expected EOBD to be '160224', got '%s'", fp.EOBD)
				}
				if fp.EOBT != "0950" {
					t.Errorf("Expected EOBD to be '160224', got '%s'", fp.EOBT)
				}
				if fp.CTOT != "1030" {
					t.Errorf("Expected CTOT to be '1030', got '%s'", fp.CTOT)
				}
				if fp.REGUL != "RMZ24M" {
					t.Errorf("Expected REGUL to be 'RMZ24M', got '%s'", fp.REGUL)
				}
				if fp.TTO.PTID != "VEULE" {
					t.Errorf("Expected TTO.PTID to be 'VEULE', got '%s'", fp.TTO.PTID)
				}
				if fp.TTO.TO != "1050" {
					t.Errorf("Expected TTO.TO to be '1050', got '%s'", fp.TTO.TO)
				}
				if fp.TTO.FL != "F300" {
					t.Errorf("Expected TTO.FL to be 'F300', got '%s'", fp.TTO.FL)
				}
				if fp.TAXITIME != "0020" {
					t.Errorf("Expected TAXITIME to be 'F300', got '%s'", fp.TTO.FL)
				}
				if fp.REGCAUSE != "CE 81" {
					t.Errorf("Expected REGCAUSE to be 'CE 81', got '%s'", fp.REGCAUSE)
				}

				// Add more checks based on the expected content of TFD message
			},
		},
		{
			/*
				Test all fields according to the ATFCM USERS MANUAL v.22.1
			*/
			name:     "SRM message",
			filename: "SRM.txt",
			expected: func(t *testing.T, fp *goflightplan.Flightplan) {
				if fp.TITLE != "SRM" {
					t.Errorf("Expected TITLE to be TFD, got %s", fp.TITLE)
				}
				if fp.ARCID != "AMC101" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%s'", fp.ARCID)
				}
				if fp.ADEP != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%s'", fp.ADEP)
				}
				if fp.ADES != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%s'", fp.ADES)
				}
				if fp.EOBD != "160224" {
					t.Errorf("Expected EOBD to be '160224', got '%s'", fp.EOBD)
				}
				if fp.EOBT != "0950" {
					t.Errorf("Expected EOBD to be '160224', got '%s'", fp.EOBT)
				}
				if fp.NEWCTOT != "1005" {
					t.Errorf("Expected CTOT to be '1030', got '%s'", fp.CTOT)
				}
				if fp.REGUL != "RMZ24M" {
					t.Errorf("Expected REGUL to be 'RMZ24M', got '%s'", fp.REGUL)
				}
				if fp.TTO.PTID != "VEULE" {
					t.Errorf("Expected TTO.PTID to be 'VEULE', got '%s'", fp.TTO.PTID)
				}
				if fp.TTO.TO != "1025" {
					t.Errorf("Expected TTO.TO to be '1025', got '%s'", fp.TTO.TO)
				}
				if fp.TTO.FL != "F300" {
					t.Errorf("Expected TTO.FL to be 'F300', got '%s'", fp.TTO.FL)
				}
				if fp.TAXITIME != "0020" {
					t.Errorf("Expected TAXITIME to be '0020', got '%s'", fp.TTO.FL)
				}
				if fp.REGCAUSE != "CE 81" {
					t.Errorf("Expected REGCAUSE to be 'CE 81', got '%s'", fp.REGCAUSE)
				}

				// Add more checks based on the expected content of TFD message
			},
		},
		{
			/*
				Test all fields according to the ATFCM USERS MANUAL v.22.1
			*/
			name:     "SLC message",
			filename: "SLC.txt",
			expected: func(t *testing.T, fp *goflightplan.Flightplan) {
				if fp.TITLE != "SLC" {
					t.Errorf("Expected TITLE to be TFD, got %s", fp.TITLE)
				}
				if fp.ARCID != "AMC101" {
					t.Errorf("Expected ARCID to be 'DLH151', got '%s'", fp.ARCID)
				}
				if fp.ADEP != "EGLL" {
					t.Errorf("Expected ADEP to be 'EGLL', got '%s'", fp.ADEP)
				}
				if fp.ADES != "LMML" {
					t.Errorf("Expected ADES to be 'LMML', got '%s'", fp.ADES)
				}
				if fp.EOBD != "080901" {
					t.Errorf("Expected EOBD to be '080901', got '%s'", fp.EOBD)
				}
				if fp.EOBT != "0945" {
					t.Errorf("Expected EOBD to be '0945', got '%s'", fp.EOBT)
				}
				if fp.REASON != "OUTREG" {
					t.Errorf("Expected REASON to be 'OUTREG', got '%s'", fp.REASON)
				}
				if fp.TAXITIME != "0020" {
					t.Errorf("Expected TAXITIME to be '0020', got '%s'", fp.TTO.FL)
				}
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
