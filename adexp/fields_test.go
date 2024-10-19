package adexp

import "testing"

func Test_MessageSetFromJSON(t *testing.T) {
	_, err := MessageSetFromJSON("../test/schema", "test")
	if err != nil {
		t.Errorf("error loading json: %v\n", err)
	}
}

func Test_MessageSetFromJSON_ERROR(t *testing.T) {
	_, err := MessageSetFromJSON("./doesnotexist", "test")
	if err.Error() != "open ./doesnotexist: no such file or directory" {
		t.Errorf("Expected error loading json: %v\n", err.Error())
	}

	_, err = MessageSetFromJSON("./", "test")
	if err.Error() != "length of set is 0" {
		t.Errorf("Expected error loading json: %v\n", err.Error())
	}
}
