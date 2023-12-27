package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

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
	store.Lock()
	// make the put operation idempotent
	store.m[key] = value
	store.Unlock()
	return nil
}

var (
	ErrKeyNotFound = errors.New("key not found")
)

func Get(key string) (string, error) {
	store.Lock()
	v, ok := store.m[key]
	store.Unlock()
	if !ok {
		return "", ErrKeyNotFound
	}
	return v, nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	return nil
}
