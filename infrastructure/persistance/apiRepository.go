package persistance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nwehr/npass/core/domain"
)

type ApiRepository struct {
	Url       string `json:"url"`
	AuthToken string `json:"authToken"`
}

func (r ApiRepository) AddEntry(entry domain.Entry) error {
	password := strings.NewReader(entry.Password)

	req, err := http.NewRequest("POST", r.Url+"/add/"+entry.Name, password)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+r.AuthToken)

	_, err = http.DefaultClient.Do(req)
	return err
}

func (r ApiRepository) RemoveEntry(name string) error {
	req, err := http.NewRequest("POST", r.Url+"/rm/"+name, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+r.AuthToken)

	_, err = http.DefaultClient.Do(req)
	return err
}

func (r ApiRepository) GetEntry(name string) (domain.Entry, error) {
	req, err := http.NewRequest("GET", r.Url+"/"+name, nil)
	if err != nil {
		return domain.Entry{}, err
	}

	req.Header.Add("Authorization", "Bearer "+r.AuthToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.Entry{}, err
	}

	password, err := ioutil.ReadAll(res.Body)

	return domain.Entry{Name: name, Password: string(password)}, err
}

func (r ApiRepository) GetEntryNames() ([]string, error) {
	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+r.AuthToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	names := []string{}
	err = json.NewDecoder(res.Body).Decode(&names)

	return names, err
}
