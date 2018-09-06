package fsutils

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestCopy(t *testing.T) {
	src := "etc/textfile.txt"
	dst := "/tmp/newtextfile.txt"

	err := Copy(src, dst)
	if err != nil {
		t.Errorf("Copy src (%s) to dst (%s) err %v ", src, dst, err)
	}

	srcb, err := ioutil.ReadFile(src)
	if err != nil {
		t.Errorf("Readfile src %s err %v", src, err)
	}

	dstb, err := ioutil.ReadFile(dst)
	if err != nil {
		t.Errorf("Readfile dst %s err %v", dst, err)
	}

	if !bytes.Equal(srcb, dstb) {
		t.Error("expect (dst) == (src) got not equal ")
	}
}
