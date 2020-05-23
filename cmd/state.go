package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

var lock sync.Mutex

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	// Create the file at the given path
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal our struct to something we can write to our file
	jsonContents, err := marshal(v)
	if err != nil {
		return err
	}

	// Copy the contents of struct to the file
	_, err = io.Copy(file, jsonContents)

	return err
}

// marshal is a function that marshals the object into an io.Reader.
func marshal(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Load loads the file at the given path into the given struct.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return unmarshal(file, v)
}

// Unmarshal is a function that unmarshals the data from the reader into the given value.
func unmarshal(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// Remove deletes the file at the given path
func Remove(path string) error {
	lock.Lock()
	defer lock.Unlock()

	err := os.Remove(path)

	return err
}
