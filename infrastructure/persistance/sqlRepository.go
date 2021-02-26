package persistance

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nwehr/paws/core/domain"
)

type SqlRepository struct {
	db *sql.DB
}

func (r SqlRepository) AddEntry(entry domain.Entry) error {
	_, err := r.db.Exec("insert into entries (name, password) values ($1, $2)", entry.Name, entry.Password)
	return err
}

func (r SqlRepository) RemoveEntry(name string) error {
	_, err := r.db.Exec("delete from entries where name = $1", name)
	return err
}

func (r SqlRepository) GetEntry(name string) (domain.Entry, error) {
	entry := domain.Entry{}
	err := r.db.QueryRow("select name, password from entries where name = $1", name).Scan(&entry.Name, &entry.Password)

	return entry, err
}

func (r SqlRepository) GetEntryNames() ([]string, error) {
	rows, err := r.db.Query("select name from entries")
	if err != nil {
		return nil, err
	}

	names := []string{}

	for {
		if !rows.Next() {
			break
		}

		name := ""
		rows.Scan(&name)

		names = append(names, name)
	}

	return names, nil
}

func NewSqlRepository(driver, dsn string) (SqlRepository, error) {
	db, err := sql.Open(driver, dsn)
	return SqlRepository{db}, err
}
