package fsutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ReadJSON from disk, then unmarshal into object.  This will likely
// be wrapped by Storage.JSON() for caching
func ReadJSON(path string, obj interface{}) error {

	// Read whole file []byte buffer
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ReadJSON expect content from (%s) got error (%v)", path, err)
	}

	// Unravel JSON into a Go thing of some sort
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		return fmt.Errorf("ReadJSON Umarshal path (%s) got error (%v)", path, err)
	}
	return nil
}

// WriteJSON - Turn the object into json then save it to the file ...
func WriteJSON(path string, obj interface{}) error {
	// JSONify
	jbytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("SaveJSON failed (%s) error (%v)", path, err)
	}
	err = ioutil.WriteFile(path, jbytes, 0755)
	if err != nil {
		return fmt.Errorf("SaveJSON write file failed (%s) error (%v)", path, err)
	}
	return nil
}
