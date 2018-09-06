package fsutils

import (
	"testing"
)

/*
Test some straight File System utilities.
*/

type testStringBool struct {
	name string
	want bool
}

var testDirs = []testStringBool{
	{"/etc/passwd", true},         // true and readable
	{"/tmp", true},                // true but a directory
	{"/srv/magfs/entries/", true}, // typically what we use
	{"/wrinkle-fritz", false},
	{"/var/audit", true},          // true but not readable
	{"/var/audit/vault", false},   // dir parent not readable and does not exists
	{"/var/audit/current", false}, // dir parent not readable and does exist
}

// TestFileExits ensures that our function always works
func TestFileExists(t *testing.T) {
	for _, tst := range testDirs {
		got := FileExists(tst.name)
		if tst.want != got {
			t.Errorf(" does (%s) exist? expected (%t) got (%t) ",
				tst.name, tst.want, got)
		}
	}
}

// TestFileNotExists ensures that our TestFileNot exists works at all times
func TestFileNotExists(t *testing.T) {
	for _, tt := range testDirs {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = !tt.want // flip the want value for not exists
			if got := FileNotExists(tt.name); got != tt.want {
				t.Errorf(" FileNotExists = %v, want %v", got, tt.want)
			}
		})
	}
}
