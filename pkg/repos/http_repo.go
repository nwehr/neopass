package repos

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nwehr/neopass"
)

type httpRepo struct {
	BaseURL string
}

func NewHTTPRepo(baseURL string) (neopass.EntryRepo, error) {
	return httpRepo{BaseURL: baseURL}, nil
}

func (r httpRepo) AddEntry(entry neopass.Entry) error {
	encoded, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", r.BaseURL, bytes.NewReader(encoded))
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (r httpRepo) RemoveEntryByName(name string) error {
	req, err := http.NewRequest("DELETE", r.BaseURL+"&name="+name, nil)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (r httpRepo) GetEntryByName(name string) (neopass.Entry, error) {
	entry := neopass.Entry{}

	req, err := http.NewRequest("GET", r.BaseURL+"&name="+name, nil)
	if err != nil {
		return entry, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entry, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&entry)

	return entry, err
}

func (r httpRepo) ListEntryNames() ([]string, error) {
	req, err := http.NewRequest("GET", r.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	names := []string{}

	json.NewDecoder(resp.Body).Decode(&names)

	return names, nil
}
