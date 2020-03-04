package main

import (
	"encoding/csv"
	"log"
	"os"
	"testing"
)

func Test_person(t *testing.T) {
	config, err := LoadConfig("./config.json")
	if err != nil {
		log.Fatalln("error loading configuration file:", err)
	}
	w := csv.NewWriter(os.Stdout)
	for i := 0; i < 20; i++ {
		p := NewPerson(config)
		record := p.toStrings()
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		for _, v := range p.visits {
			record = v.toStrings()
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
