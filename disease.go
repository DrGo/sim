package main

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds info on run config
type Config struct {
	Diseases []*Disease `json:"diseases"`
	Version  string     `json:"version"`
	N        int        `json:"n"`
}

// Disease holds config for disease
type Disease struct {
	Name             string  `json:"name"`
	PrevalenceMale   float64 `json:"prevalence_male"`
	PrevalenceFemale float64 `json:"prevalence_female"`
	Recurrence       int     `json:"recurrence"`
	HospitalRate     Stats   `json:"hospital_rate"`
	ClinicRate       Stats   `json:"clinic_rate"`
	Icd9             string  `json:"icd9"`
	Icd10            string  `json:"icd10"`
	incidenceDate    time.Time
}

type Stats struct {
	Mean float64
	SD   float64
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	return config, err
}
