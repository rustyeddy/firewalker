package fsutils

import (
	"os"
	"testing"
)

// TestMkdir command
func TestMkdir(t *testing.T) {

	foobar := "/tmp/foo/bar" // remove dir and start all over

	tests := []struct {
		path    string
		wantErr bool
	}{
		{foobar, false},    // it did not exist
		{foobar, false},    // will not error if directory already exists
		{"/badpath", true}, // fail if we try to create a dir with a bad path
	}
	os.RemoveAll(foobar)
	// Test range
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if err := Mkdir(tt.path); (err != nil) != tt.wantErr {
				t.Errorf("Mkdir() %s error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}
