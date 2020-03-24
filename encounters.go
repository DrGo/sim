package main

import (
	"math/rand"
	"strconv"
)

type Visit struct {
	id        int64
	startDate int64
	endDate   int64
	diagnosis string
}

func (v *Visit) toStrings() []string {
	a := []string{}
	a = append(a, strconv.Itoa(int(v.id)))
	a = append(a, toTime(v.startDate).Format(dateLayoutISO))
	a = append(a, toTime(v.endDate).Format(dateLayoutISO))
	a = append(a, v.diagnosis)
	return a
}

func (p *Person) newVisit(kind int, disease *Disease, incidenceDate int64) *Visit {
	v := Visit{
		id:        p.id,
		startDate: RangeDate(incidenceDate, p.cancelDate),
	}
	switch kind {
	case kindHospital:
		v.endDate = v.startDate + int64(Normal(disease.Hospitalization.StayLength.Mean, disease.Hospitalization.StayLength.SD))*secondsInDay
		v.diagnosis = disease.Icd10
	default:
		// v.endDate = stataMissingInt64 //default to missing
		v.diagnosis = disease.Icd9
	}
	return &v
}

type Rx struct {
	id   int64
	date int64
	din  string
}

func (p *Person) newRx(disease *Disease, incidenceDate int64) *Rx {
	r := Rx{
		id:   p.id,
		date: RangeDate(incidenceDate, p.cancelDate),
	}
	for _, din := range disease.Dins {
		if rand.Float64() < din.Prob {
			r.din = din.DIN
		}
	}
	return &r
}

func (r *Rx) toStrings() []string {
	a := []string{}
	a = append(a, strconv.Itoa(int(r.id)))
	a = append(a, toTime(r.date).Format(dateLayoutISO))
	return a
}
