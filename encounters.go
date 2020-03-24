package main

import (
	"math/rand"
	"strconv"
)

type Visit struct {
	kind      int
	id        int64
	startDate int64
	endDate   int64
	diagnosis string
}

func (v *Visit) toStrings() []string {
	a := []string{}
	a = append(a, strconv.Itoa(int(v.id)))
	a = append(a, toTime(v.startDate).Format(dateLayoutISO))
	if v.kind == kindHospital {
		a = append(a, toTime(v.endDate).Format(dateLayoutISO))
	}
	a = append(a, v.diagnosis)
	return a
}

func (p *Person) newVisit(kind int, disease *Disease, incidenceDate int64) *Visit {
	v := Visit{
		kind:      kind,
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
	Drugs []*Drug
}

type Drug struct {
	id   int64
	date int64
	din  string
}

func (p *Person) newRx(disease *Disease, incidenceDate int64) *Rx {
	date := RangeDate(incidenceDate, p.cancelDate)
	var r Rx
	for _, din := range disease.Dins {
		if rand.Float64() < din.Prob {
			r.Drugs = append(r.Drugs, &Drug{
				id:   p.id,
				date: date,
				din:  din.DIN,
			})
		}
	}
	return &r
}

func (d *Drug) toStrings() []string {
	a := []string{}
	a = append(a, strconv.Itoa(int(d.id)))
	a = append(a, toTime(d.date).Format(dateLayoutISO))
	a = append(a, d.din)
	return a
}
