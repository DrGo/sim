package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	config     *Config
	bufferSize int
	lastID     int64
	wg         sync.WaitGroup
	personCh   chan []string
	hospCh     chan []string
	clinicCh   chan []string
	rxCh       chan []string
}

func NewDispatcher(bufferSize int, config *Config) *Dispatcher {
	return &Dispatcher{
		config:     config,
		bufferSize: bufferSize,
		lastID:     1000_000,
		personCh:   make(chan []string, bufferSize),
		hospCh:     make(chan []string, bufferSize),
		clinicCh:   make(chan []string, bufferSize),
		rxCh:       make(chan []string, bufferSize),
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

func (d *Dispatcher) SaveRx(records []string) {
	d.rxCh <- records
}

func (d *Dispatcher) getLastID() int64 {
	return atomic.AddInt64(&d.lastID, 1)
}

func (d *Dispatcher) getQbyId(category string) (chan []string, error) {
	switch category {
	case "person":
		return d.personCh, nil
	case "hosp":
		return d.hospCh, nil
	case "clinic":
		return d.clinicCh, nil
	case "rx":
		return d.rxCh, nil
	default:
		return nil, fmt.Errorf("no such output category: %s", category)
	}
}

func (d *Dispatcher) closeAll() {
	close(d.personCh)
	close(d.hospCh)
	close(d.clinicCh)
	close(d.rxCh)
}
