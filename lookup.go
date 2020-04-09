package main

import (
	"encoding/csv"
	"fmt"
	"github/drgo/alias"
	"io"
	"os"
	"strconv"
	"strings"
)

type Lookup struct {
	FieldName string
	Codes     []string
	Class     []string
	Probs     []float64
	alias     *alias.Alias
}

func LoadLookup(fileName, fieldName string, mustClass bool) (*Lookup, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lookup := &Lookup{FieldName: fieldName}
	csv := csv.NewReader(file)
	// Lines beginning with "/" without preceding whitespace are ignored.
	csv.Comment = '/'
	csv.FieldsPerRecord = 2
	if mustClass {
		csv.FieldsPerRecord = 3
	}
	csv.ReuseRecord = true // for performance
	hasClass := false
	if hasClass, err = validateHeader(csv); err != nil {
		return nil, err
	}
	if mustClass && !hasClass {
		return nil, fmt.Errorf("class field is missing")
	}
	recNum := 1
readLoop:
	for {
		record, err := csv.Read()
		switch {
		case err == io.EOF:
			break readLoop
		case err != nil:
			return nil, err
		}
		recNum++
		if strings.TrimSpace(record[0]) == "" {
			return nil, fmt.Errorf("missing code in line number %d", recNum)
		}
		lookup.Codes = append(lookup.Codes, record[0])
		if strings.TrimSpace(record[1]) == "" {
			return nil, fmt.Errorf("missing prob in line number %d", recNum)
		}
		prob, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid probability in line number %d: %s", recNum, err)
		}
		lookup.Probs = append(lookup.Probs, prob)
		if hasClass {
			lookup.Class = append(lookup.Class, record[2])
		}
	}
	lookup.alias, err = alias.New(lookup.Probs)
	if err != nil {
		return nil, err
	}
	return lookup, nil
}

// RandCode returns a randomly selected code
func (l *Lookup) RandCode() string {
	return l.Codes[l.alias.Draw()]
}

func validateHeader(csv *csv.Reader) (bool, error) {
	record, err := csv.Read()
	switch {
	case err == io.EOF:
		return false, fmt.Errorf("empty csv file")
	case err != nil:
		return false, err
	}
	if record[0] != "code" && record[1] != "prob" {
		return false, fmt.Errorf("required field names are missing. The first two fields must be named 'code' and 'prob'")
	}
	if len(record) > 2 && record[2] == "class" {
		return true, nil
	}
	return false, nil
}
