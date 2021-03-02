package klaabu

import (
	"errors"
	"fmt"
	"io/ioutil"
)

// Schema is the Klaabu schema.
type Schema struct {
	Version string
	Labels  map[string]string
	Root    *Prefix
}

// NewSchema creates and returns a new Schema.
func NewSchema(labels map[string]string) *Schema {
	return &Schema{
		Version: "v1",
		Labels:  labels,
	}
}

func (s *Schema) PrefixById(id string) *Prefix {
	return s.Root.PrefixById(id)
}

// Validate checks if you are stupid or not.
func (s *Schema) Validate() error {
	err := s.Root.Validate()
	if err != nil {
		return err
	}

	return nil
}

// WriteSchemaToFile takes a schema and writes it to a file on disk.
func WriteSchemaToFile(s *Schema, fileName *string) error {
	// TODO use kml
	data, err := []byte{}, errors.New("not implemented")
	if err != nil {
		return fmt.Errorf("error in serialization: %s", err)
	}

	err = ioutil.WriteFile(*fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing schema to file: %s", err)
	}

	return nil
}
