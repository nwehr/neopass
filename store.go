package neopass

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Entries []Entry `yaml:"entries"`
}

func (s *Store) ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	defer file.Close()

	return s.Read(file)
}

func (s *Store) Read(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(s)
}

func (s Store) WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	return s.Write(file)
}

func (s Store) Write(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(s)
}
