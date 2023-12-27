package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var store map[string]string

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	key := v["key"]
	value, err := io.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	key := v["key"]

	value, err := Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(value))
}

func main() {
	store = make(map[string]string)
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}

func Put(key, value string) error {
	// make the put operation idempotent
	store[key] = value
	return nil
}

var (
	ErrKeyNotFound = errors.New("key not found")
)

func Get(key string) (string, error) {
	if v, ok := store[key]; ok {
		return v, nil
	}
	return "", ErrKeyNotFound
}

func Delete(key string) error {
	delete(store, key)
	return nil
}
