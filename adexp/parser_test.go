package adexp

import (
	"os"
	"path/filepath"
	"reflect"
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
				// Add more checks based on the expected content of BFD message
				if fp.Flightplan.ADEP != "LATI" {
					t.Errorf("Expected ADEP to be LATI, got %s", fp.Flightplan.ADEP)
				}
				if fp.Flightplan.ADES != "EDJA" {
					t.Errorf("Expected ADES to be EDJA, got %s", fp.Flightplan.ADES)
				}
				// Check other fields specific to BFD message
			},
		},
		{
			name:     "CFD message",
			filename: "008_CFD.txt",
			expected: func(t *testing.T, fp *goflightplan.FlightplanWrapper) {
				if fp.Flightplan.TITLE != "CFD" {
					t.Errorf("Expected TITLE to be CFD, got %s", fp.Flightplan.TITLE)
				}
				if fp.Flightplan.ARCID != "WMT3GH" {
					t.Errorf("Expected ARCID to be 'WMT3GH', got '%s'", fp.Flightplan.ARCID)
				}
				// Add more checks based on the expected content of CFD message
				// For example, check for changes in flight data
			},
		},
		{
			name:     "TFD message",
			filename: "009_TFD.txt",
			expected: func(t *testing.T, fp *goflightplan.FlightplanWrapper) {
				if fp.Flightplan.TITLE != "TFD" {
					t.Errorf("Expected TITLE to be TFD, got %s", fp.Flightplan.TITLE)
				}
				if fp.Flightplan.ARCID != "WZZ70BK" {
					t.Errorf("Expected ARCID to be 'WZZ70BK', got '%s'", fp.Flightplan.ARCID)
				}
				// Add more checks based on the expected content of TFD message
				// For example, check for termination-related fields
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

func TestParser_findFieldInSchema(t *testing.T) {
	testSchema := []MessageSet{
		{
			Name: "TestSet",
			Set: map[string]StandardSchema{
				"TEST": {
					Category: "TEST",
					Items: []DataField{
						{DataItem: "TITLE", Type: Basicfield},
						{DataItem: "ARCID", Type: Basicfield},
					},
				},
			},
		},
	}

	parser := NewParser(testSchema)

	testCases := []struct {
		name      string
		fieldName string
		want      DataField
		wantErr   bool
	}{
		{
			name:      "Existing field",
			fieldName: "TITLE",
			want:      DataField{DataItem: "TITLE", Type: Basicfield},
			wantErr:   false,
		},
		{
			name:      "Non-existing field",
			fieldName: "NONEXISTENT",
			want:      DataField{},
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.findFieldInSchema(tc.fieldName)
			if tc.wantErr {
				if err == nil {
					t.Errorf("findFieldInSchema() error = nil, wantErr %v", tc.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("findFieldInSchema() unexpected error = %v", err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("findFieldInSchema() got = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestParser_setFieldValue(t *testing.T) {
	parser := NewParser(nil) // Schema not needed for this test
	parser.flightplan = &goflightplan.Flightplan{}

	testCases := []struct {
		name    string
		field   DataField
		value   string
		wantErr bool
	}{
		{
			name:    "Set string field",
			field:   DataField{DataItem: "TITLE", Type: Basicfield},
			value:   "TEST",
			wantErr: false,
		},
		{
			name:    "Set non-existent field",
			field:   DataField{DataItem: "NONEXISTENT", Type: Basicfield},
			value:   "TEST",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := parser.setFieldValue(tc.field, tc.value)
			if tc.wantErr {
				if err == nil {
					t.Errorf("setFieldValue() error = nil, wantErr %v", tc.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("setFieldValue() unexpected error = %v", err)
				}
				got := reflect.ValueOf(parser.flightplan).Elem().FieldByName(tc.field.DataItem).String()
				if got != tc.value {
					t.Errorf("setFieldValue() got = %+v, want %+v", got, tc.value)
				}
			}
		})
	}
}

func loadTestCases(directory string) ([]string, error) {
	var testCases []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			testCases = append(testCases, string(content))
		}
		return nil
	})
	return testCases, err
}

// Helper function to load the MessageSet (as defined in previous response)
func loadTestMessageSet(t *testing.T) MessageSet {
	t.Helper()
	messageSet, err := MessageSetFromJSON("../test/json", "TestSet")
	if err != nil {
		t.Fatalf("Failed to load message set from JSON: %v", err)
	}
	return *messageSet
}
