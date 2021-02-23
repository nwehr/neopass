package core

import "fmt"

type Entries []Entry

func (e *Entries) Add(entry Entry) error {
	for _, current := range *e {
		if current.Name == entry.Name {
			return AlreadyExistsError{entry.Name}
		}
	}

	*e = append(*e, entry)
	return nil
}

func (e *Entries) Update(entry Entry) error {
	for i, current := range *e {
		if current.Name == entry.Name {
			(*e)[i].Password = entry.Password
			return nil
		}
	}

	return NotFoundError{entry.Name}
}

type AlreadyExistsError struct {
	Name string
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("entry %s already exists", e.Name)
}

type NotFoundError struct {
	Name string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("entry %s does not exists", e.Name)
}
