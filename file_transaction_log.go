package main

import (
	"os"
	"sync"
)

// FileTransactionLogger defines methods and fields to write transaction log to a file in the disk
type FileTransactionLogger struct {
	file   *os.File
	events chan<- Event
	errors <-chan error
	sync.RWMutex
	lastSequenceNumber uint64
}

func (f *FileTransactionLogger) WritePutEvent(key, value string) {
	f.events <- Event{
		eventType: PutEvent,
		key:       key,
		value:     value,
	}
}
func (f *FileTransactionLogger) WriteDeleteEvent(key string) {
	f.events <- Event{
		eventType: DeleteEvent,
		key:       key,
	}
}

func (f *FileTransactionLogger) Err() <-chan error {
	return f.errors
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileTransactionLogger{
		file: file,
	}, nil

}
