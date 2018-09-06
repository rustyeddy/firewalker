package fsutils

import (
	"testing"
)

type testObject struct {
	Name        string
	Description string
	Count       int
}

func JSONTest(t *testing.T) {

	tfile := "/tmp/testobj.json"
	to := &testObject{"TestObject", "This is a Test", 6}
	err := WriteJSON(tfile, to)
	if err != nil {
		t.Errorf("expected to write %s err %v ", tfile, err)
	}

	var in *testObject
	err = ReadJSON(tfile, in)
	if err != nil {
		t.Errorf("expected to read %s err %v ", tfile, err)
	}

	if in == nil {
		t.Errorf("expected incoming object got nil ")
	}

	if in.Name != to.Name ||
		in.Description != to.Description ||
		in.Count != to.Count {
		t.Errorf("expected (%+v)", to)
		t.Errorf("     got (%+v)", in)
	}
}
