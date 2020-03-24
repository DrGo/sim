package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	configFileName = "./config.json"
	bufferSize     = 100
)

func main() {
	done := make(chan struct{}) //main receives done signal on this chan
	config, err := LoadConfig(configFileName)
	if err != nil {
		log.Fatalln("error loading configuration file:", err)
	}
	dispatcher := NewDispatcher(bufferSize)
	go writer("person", dispatcher, done)
	go writer("hosp", dispatcher, done)
	go writer("clinic", dispatcher, done)
	go writer("rx", dispatcher, done)
	for i := 0; i < config.N; i++ {
		dispatcher.wg.Add(1)
		go NewPerson(config, dispatcher)
	}
	dispatcher.wg.Wait()
	close(dispatcher.personCh)
	close(dispatcher.hospCh)
	close(dispatcher.clinicCh)
	close(dispatcher.rxCh)

	for i := 0; i < 4; i++ {
		<-done //wait for all writers to quit
	}
}

func writer(category string, dispatcher *Dispatcher, done chan struct{}) {
	f, err := os.Create(category + ".csv")
	if err != nil {
		log.Fatalln("error writing to file:", err)
	}
	defer f.Close()
	log.Println("creating file:", category)
	var qu chan []string
	var fieldNames []string
	switch category {
	case "person":
		qu = dispatcher.personCh
		fieldNames = strings.Split("subject_id,gender,birthdate,age,coverage_start,coverage_end", ",")
	case "hosp":
		qu = dispatcher.hospCh
		fieldNames = strings.Split("subject_id,service_date,discharge_date,code", ",")
	case "clinic":
		qu = dispatcher.clinicCh
		fieldNames = strings.Split("subject_id,service_date,discharge_date,code", ",")
	case "rx":
		qu = dispatcher.rxCh
		fieldNames = strings.Split("subject_id,service_date,code", ",")
	default:
		log.Fatalln("no such output category: ", category)
	}
	w := csv.NewWriter(f)
	// write fieldNames
	if err := w.Write(fieldNames); err != nil {
		log.Fatalln("error writing record to csv: ", err)
	}
	for record := range qu {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv: ", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done writing to:", category)
	done <- struct{}{}
}
