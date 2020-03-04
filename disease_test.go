package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		args    string
		want    *Config
		wantErr bool
	}{
		{"./config.json", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			fmt.Println(tt.args)
			got, err := LoadConfig(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			json, err := json.MarshalIndent(got, "", " ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
		})
	}
}

func TestConfigMarshal(t *testing.T) {
	config := Config{Diseases: []*Disease{
		&Disease{
			Name:         "diabetes",
			HospitalRate: Stats{3, 1},
		},
	}}
	json, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}
