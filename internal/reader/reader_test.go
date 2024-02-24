package reader

import (
	"os"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	r := NewReader(nil)

	tests := []struct {
		description   string
		fileContent   string
		expectedLines []string
	}{
		{
			description:   "one line file",
			fileContent:   "20:00 InitGame: line 0",
			expectedLines: []string{"20:00 InitGame: line 0"},
		},
		{
			description:   "two line file",
			fileContent:   "20:00 InitGame: line 0\n20:00 InitGame: line 1",
			expectedLines: []string{"20:00 InitGame: line 0", "20:00 InitGame: line 1"},
		},
	}

	for _, test := range tests {
		tmpfile, err := os.CreateTemp("", "test_file")
		if err != nil {
			t.Errorf("%s: Error creating temp file: %v", test.description, err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(test.fileContent)); err != nil {
			t.Errorf("%s: Error writing to temp file: %v", test.description, err)
		}
		defer tmpfile.Close()

		lines, err := r.Read(tmpfile.Name())
		if err != nil {
			t.Errorf("%s: Error reading file: %v", test.description, err)
		}
		if !reflect.DeepEqual(lines, test.expectedLines) {
			t.Errorf("%s: Expected %v, got %v", test.description, test.expectedLines, lines)
		}
	}
}
