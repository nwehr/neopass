package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/nwehr/neopass"
)

var (
	commit    string
	buildDate string
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		Fatalf("%v\n", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/{client_uuid}/{name}", getEntryHandler(conn)).Methods("GET")
	r.HandleFunc("/{client_uuid}/{name}", deleteEntryHandler(conn)).Methods("DELETE")
	r.HandleFunc("/{client_uuid}", postEntryHandler(conn)).Methods("POST")
	r.HandleFunc("/{client_uuid}", listNamesHandler(conn)).Methods("GET")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		Fatalf("%v\n", err)
	}
}

func listNamesHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientUUID := mux.Vars(r)["client_uuid"]

		rows, err := conn.Query(context.Background(), `select "name" from entries where "client_uuid" = $1`, clientUUID)
		if err != nil {
			// return nil, fmt.Errorf("could not query entries: %v", err)
		}

		defer rows.Close()

		names := []string{}

		for rows.Next() {
			name := ""

			err = rows.Scan(&name)
			if err != nil {
				// return names, fmt.Errorf("could not scan into name: %v", err)
			}

			names = append(names, name)
		}

		json.NewEncoder(w).Encode(names)
	}
}

func getEntryHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientUUID := mux.Vars(r)["client_uuid"]

		entry := neopass.Entry{
			Name: mux.Vars(r)["name"],
		}

		conn.QueryRow(context.Background(), `select "name", "password" from entries where "client_uuid" = $1 and "name" = $2`, clientUUID, entry.Name).Scan(&entry.Password)

		json.NewEncoder(w).Encode(entry)
	}
}

func postEntryHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientUUID := mux.Vars(r)["client_uuid"]

		entry := neopass.Entry{}

		json.NewDecoder(r.Body).Decode(&entry)

		conn.Exec(context.Background(), `insert into entries ("client_uuid", "name", "password") values ($1, $2, $3)`, clientUUID, entry.Name, entry.Password)

		json.NewEncoder(w).Encode(entry)
	}
}

func deleteEntryHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientUUID := mux.Vars(r)["client_uuid"]
		name := mux.Vars(r)["name"]

		conn.Exec(context.Background(), `delete from entries where "client_uuid" = $1 and "name" = $2`, clientUUID, name)
	}
}

func Fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
