package main

type EventType byte

const (
	_        = iota
	PutEvent = iota
	DeleteEvent
)

type Event struct {
	sequenceNumber uint64
	eventType      EventType
	key            string
	value          string
}

type TransactionLogger interface {
	WritePutEvent(key, value string)
	WriteDeleteEvent(key string)

	Err() <-chan error
}
