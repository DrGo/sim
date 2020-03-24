package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds info on run config
type Config struct {
	Diseases        []*Disease       `json:"diseases"`
	Version         string           `json:"version"`
	N               int              `json:"n"`
	Population      *Population      `json:"population"`
	Hospitalization *Hospitalization `json:"hospitalization"`
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
	incidenceDate    int64
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
	StayLength struct {
		Mean float64 `json:"Mean"`
		SD   float64 `json:"SD"`
	} `json:"stay_length"`
}

type Stats struct {
	Mean float64
	SD   float64
}

type DIN struct {
	Prob float64
	DIN  string
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
	return config, nil
}
