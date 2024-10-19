package adexp

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Parse_ListField(t *testing.T) {
	// Load the test schema from JSON files
	testSchema := []MessageSet{LoadTestMessageSet(t)}
	parser := NewParser(testSchema)

	testCases := []struct {
		name        string
		filename    string
		description string
		expected    func(*testing.T, map[string]interface{})
	}{
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

				//
				rtepts, ok := fp["RTEPTS"].([]interface{})
				if !ok {
					t.Errorf("Expected RTEPTS to be a []interface{}, got %T", fp["RTEPTS"])
				} else {
					if len(rtepts) != 3 {
						t.Errorf("Expected 3 route points, got %d", len(rtepts))
					}

					expectedPoints := []map[string]string{
						{"PTID": "WOODY", "TO": "1235", "FL": "F210"},
						{"PTID": "CIV", "TO": "1239", "FL": "F330"},
						{"PTID": "NEBUL", "TO": "1240", "FL": "F330"},
					}

					for i, expectedPt := range expectedPoints {
						pt, ok := rtepts[i].(map[string]interface{})
						if !ok {
							t.Errorf("Expected route point %d to be a map[string]interface{}, got %T", i, rtepts[i])
							continue
						}

						for key, expectedValue := range expectedPt {
							if pt[key] != expectedValue {
								t.Errorf("Expected route point %d %s to be %s, got %v", i, key, expectedValue, pt[key])
							}
						}
					}
				}

				//EQP-BEGIN EQCST

				eqcst, ok := fp["EQCST"].([]interface{})
				if !ok {
					t.Errorf("Expected EQCST to be a []interface{}, got %T", fp["EQCST"])
				} else {
					if len(eqcst) != 2 {
						t.Errorf("Expected 2 equipments, got %d", len(eqcst))
					}

					expectedPoints := []string{
						"W/EQ", "Y/NO",
					}

					for i, expectedPt := range expectedPoints {
						pt, ok := eqcst[i].(string)
						if !ok {
							t.Errorf("Expected route point %d to be a string, got %T", i, eqcst[i])
							continue
						}
						if pt != expectedPt {
							t.Errorf("Expected route point %d to be %s, got %v", i, expectedPt, pt)
						}

					}
				}
			}},
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
