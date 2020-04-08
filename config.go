package main

import (
	"encoding/json"
	"fmt"
	"github/drgo/alias"
	"os"
	"strings"
	"time"
)

// Config holds info on run config
type Config struct {
	Version         string           `json:"version"`
	Seed            int              `json:"seed"`
	N               int              `json:"n"`
	Diseases        []*Disease       `json:"diseases"`
	Population      *Population      `json:"population"`
	Hospitalization *Hospitalization `json:"hospitalization"`
	Locator         *Locator         `json:"locator"`
	Options         struct {
		LocationNeeded bool `json:"location_needed"`
	} `json:"options"`
	locatorName  string
	locatorCodes []string
	locatorFreqs []int
	locatorAlias *alias.Alias
	fieldNames   map[string]string
}

// Disease holds config for disease
type Disease struct {
	Name             string           `json:"name"`
	PrevalenceMale   float64          `json:"prevalence_male"`
	PrevalenceFemale float64          `json:"prevalence_female"`
	Recurrence       int              `json:"recurrence"`
	HospitalRate     Stats            `json:"hospital_rate"`
	ClinicRate       Stats            `json:"clinic_rate"`
	Icd9             string           `json:"icd9"`
	Icd10            string           `json:"icd10"`
	RxRate           Stats            `json:"rx_rate"`
	Dins             []DIN            `json:"dins"`
	Hospitalization  *Hospitalization `json:"hospitalization"`
}

type Population struct {
	MigrantProb       float64 `json:"migrant_prob"`
	CancelProb        float64 `json:"cancel_prob"`
	DatabaseStartDate string  `json:"database_start_date"`
	EarliestBirthDate string  `json:"earliest_birth_date"`
	minDate           int64
	databaseStartDate int64
}

type Hospitalization struct {
	StayLength Stats `json:"stay_length"`
}

type Stats struct {
	Mean float64
	SD   float64
}

type DIN struct {
	Prob float64
	DIN  string
}

type Locator struct {
	Active bool   `json:"active"`
	Name   string `json:"name"`
	Codes  []struct {
		Code string `json:"code"`
		Freq int    `json:"freq"`
	} `json:"codes"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		return nil, err
	}
	return ProcessConfig(config)
}

func ProcessConfig(config *Config) (*Config, error) {
	var (
		err  error
		date time.Time
	)
	if date, err = time.Parse("2006-01-02", config.Population.DatabaseStartDate); err != nil {
		return nil, err
	}
	config.Population.databaseStartDate = date.Unix()
	if date, err = time.Parse("2006-01-02", config.Population.EarliestBirthDate); err != nil {
		return nil, err
	}
	config.Population.minDate = date.Unix()
	if len(config.Diseases) == 0 {
		return nil, fmt.Errorf("Configuration must include at least 1 disease entry")
	}
	if config.Hospitalization == nil {
		return nil, fmt.Errorf("Configuration must include a hospitalization entry")
	}
	for _, disease := range config.Diseases {
		if disease.Hospitalization == nil {
			disease.Hospitalization = config.Hospitalization
		}
	}
	if config.Options.LocationNeeded {
		if config.Locator == nil {
			return nil, fmt.Errorf("Location_needed is set to true so Configuration must include a valid Locator entry")
		}
		if len(config.Locator.Codes) == 0 {
			return nil, fmt.Errorf("Location_needed is set to true so locator.codes entry must include at least one entry")
		}
		if strings.TrimSpace(config.Locator.Name) == "" {
			config.Locator.Name = "postal_code"
		}
		for _, e := range config.Locator.Codes {
			config.locatorCodes = append(config.locatorCodes, e.Code)
			config.locatorFreqs = append(config.locatorFreqs, e.Freq)
		}
		if config.locatorAlias, err = alias.NewFreq(config.locatorFreqs); err != nil {
			return nil, fmt.Errorf("error parsing locator information: %v", err)
		}
	}
	// define field names to use in csv
	config.fieldNames = make(map[string]string, 4)
	config.fieldNames["person"] = "subject_id,gender,birthdate,age,coverage_start,coverage_end"
	if config.Options.LocationNeeded {
		config.fieldNames["person"] += "," + config.Locator.Name
	}
	config.fieldNames["hosp"] = "subject_id,service_date,discharge_date,code"
	config.fieldNames["clinic"] = "subject_id,service_date,code"
	config.fieldNames["rx"] = "subject_id,service_date,code"
	return config, nil
}
