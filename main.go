package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	go writer("person.csv", dispatcher.personCh, dispatcher, done)
	go writer("hosp.csv", dispatcher.hospCh, dispatcher, done)
	go writer("clinic.csv", dispatcher.clinicCh, dispatcher, done)
	go writer("rx.csv", dispatcher.rxCh, dispatcher, done)
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

func writer(fileName string, qu chan []string, dispatcher *Dispatcher, done chan struct{}) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("error writing to file:", err)
	}
	defer f.Close()
	//TODO: use filename to derive ch name and field variables
	log.Println("creating file:", fileName)
	w := csv.NewWriter(f)
	for record := range qu {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done writing to:", fileName)
	done <- struct{}{}
}
