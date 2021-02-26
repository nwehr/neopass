package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/core/usecases"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type Api struct {
	ApiConfig  Config
	Repository domain.StoreRepository
	Encrypter  encryption.Encrypter
	Decrypter  encryption.Decrypter
}

func (iface Api) Start() error {
	r := mux.NewRouter()

	r.Use(RequireAuthorization(iface.ApiConfig.AuthToken))

	r.HandleFunc("/", GetAllEntriesHandler(iface.Repository)).Methods("GET")
	r.HandleFunc("/{name}", GetPasswordHandler(iface.Repository, iface.Decrypter)).Methods("GET")
	r.HandleFunc("/add/{name}", AddEntryHandler(iface.Repository, iface.Encrypter)).Methods("POST")
	r.HandleFunc("/update/{name}", UpdateEntryHandler(iface.Repository, iface.Encrypter)).Methods("POST")
	r.HandleFunc("/rm/{name}", RemoveEntryHandler(iface.Repository)).Methods("POST")

	fmt.Printf("listening on %s\n", iface.ApiConfig.Listen)

	return http.ListenAndServe(iface.ApiConfig.Listen, r)
}

func RequireAuthorization(authToken string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Method, r.URL.Path)

			bearerToken := r.Header.Get("Authorization")
			if bearerToken == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")
			if bearerToken != authToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetAllEntriesHandler(repo domain.StoreRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names, err := usecases.GetAllEntryNames{repo}.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(names)
	}
}

func GetPasswordHandler(repo domain.StoreRepository, dec encryption.Decrypter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		password, err := usecases.GetDecryptedPassword{repo, dec}.Run(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(password))

	}
}

func AddEntryHandler(repo domain.StoreRepository, enc encryption.Encrypter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		password, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err = (usecases.AddEntry{repo, enc}.Run(name, string(password))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func UpdateEntryHandler(repo domain.StoreRepository, enc encryption.Encrypter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		password, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err = (usecases.UpdateEntry{repo, enc}.Run(name, string(password))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RemoveEntryHandler(repo domain.StoreRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		if err := (usecases.RemoveEntry{repo}.Run(name)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
