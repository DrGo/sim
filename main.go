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

	go writer("person", config, done)
	go writer("hosp", config, done)
	go writer("clinic", config, done)
	go writer("rx", config, done)
	for i := 0; i < config.N; i++ {
		config.dispatcher.wg.Add(1)
		go NewPerson(config)
	}
	config.dispatcher.wg.Wait()
	config.dispatcher.closeAll()

	for i := 0; i < 4; i++ {
		<-done //wait for all writers to quit
	}
}

func writer(category string, config *Config, done chan struct{}) {
	f, err := os.Create(category + ".csv")
	if err != nil {
		log.Fatalln("error writing to file:", err)
	}
	defer f.Close()
	log.Println("creating file:", category)
	var fieldNames []string
	fieldNames = strings.Split(config.fieldNames[category], ",")
	qu, err := config.dispatcher.getQbyId(category)
	if err != nil {
		log.Fatalln(err)
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
