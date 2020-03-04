package main

import (
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	bufferSize int
	lastID     int64
	wg         sync.WaitGroup
	personCh   chan []string
	hospCh     chan []string
	clinicCh   chan []string
}

func NewDispatcher(bufferSize int) *Dispatcher {
	return &Dispatcher{
		bufferSize: bufferSize,
		lastID:     1000_000,
		personCh:   make(chan []string, bufferSize),
		hospCh:     make(chan []string, bufferSize),
		clinicCh:   make(chan []string, bufferSize),
	}
}

func (d *Dispatcher) SavePerson(records []string) {
	d.personCh <- records
}

func (d *Dispatcher) SaveHosp(records []string) {
	d.hospCh <- records
}

func (d *Dispatcher) SaveClinic(records []string) {
	d.clinicCh <- records
}

func (d *Dispatcher) getLastID() int64 {
	return atomic.AddInt64(&d.lastID, 1)
}
