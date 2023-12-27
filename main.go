package main

import "errors"

var store map[string]string

func main() {
	store = make(map[string]string)
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
