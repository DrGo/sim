# Sim
generates random  but plausible healthcare utilization data using a template stored in config.json.

## Usage
simply, type sim in a folder where config.json exists


## Rules for config.json
The "__doc" key can be used to document the configuration file.

version: must be 1.0.
n: the number of patient records to generate. Must be >0.

	"population": {
		"migrant_prob": 0.15,
		"cancel_prob": 0.15,
		"database_start_date": "1971-01-01",
		"earliest_birth_date": "1920-01-01"
	},

hospitalization: sets parameters for all hospitalizations regardless of disease
  stay_length: provides the mean and SD of the distribution of hospital length of stay in days

diseases: array of disease descriptor

chronic and recurrence are not implemented

hospital_rate: provides the mean and SD of the distribution of number of hospitalizations per year.

clinic_rate: provides the mean and SD of the distribution of number of hospitalizations per year.

rx_rate: provides the mean and SD of the distribution of the number of prescription filled per year.

dins: an array of 1 or more drugs filled. din=as per the DPD; prob= probability of getting this DIN.

locator: used to generate a random geolocation code 
		name: is the name of the field in the generated dataset, eg. postal_code.
		csv_filename: the relative/absolute path of a csv file that must"
      - start by the following string "code, freq"
      - each subsequent line contains geocode and the probability of residing in that geocode, eg R3L1E9, 0.10

