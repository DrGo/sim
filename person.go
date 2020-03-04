package main

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	dateLayoutISO = "2006-01-02"
	secondsInDay  = 24 * 60 * 60
	daysInYear    = 365
)

const (
	hospital = iota
	clinic
)

var (
	today             = time.Now()
	todayUnix         = today.Unix()
	migrantProb       = 0.15
	cancelProb        = 0.15
	minDate           = time.Date(1920, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	databaseStartDate = time.Date(1971, 1, 0, 0, 0, 0, 0, time.UTC).Unix()

	hospitalStayMean = 7.0
	hospitalStaySD   = 1.0
)

// Person generates a person data
type Person struct {
	config     *Config
	dispatcher *Dispatcher
	id         int64
	sex        int
	age        int
	dob        int64
	dod        int64
	regisDate  int64
	cancelDate int64
	visits     []*Visit
}

func NewPerson(config *Config, dispatcher *Dispatcher) *Person {
	p := Person{
		config:     config,
		dispatcher: dispatcher,
		id:         dispatcher.getLastID(),
		sex:        RangeInt(0, 1), //0 male 1 female
		dob:        RangeDate(minDate, todayUnix),
		visits:     []*Visit{},
	}
	defer dispatcher.wg.Done()
	dob := toTime(p.dob)
	p.age = today.Year() - dob.Year()
	if rand.Float64() < migrantProb {
		p.regisDate = RangeDate(databaseStartDate, todayUnix)
	} else {
		p.regisDate = databaseStartDate
	}
	if p.dob > p.regisDate {
		p.dob = p.regisDate
	}
	if rand.Float64() < cancelProb {
		p.cancelDate = RangeDate(p.regisDate, todayUnix)
	} else {
		p.cancelDate = todayUnix
	}
	p.dispatcher.SavePerson(p.toStrings())
	p.addVisits()
	return &p
}

func (p *Person) toStrings() []string {
	a := []string{}
	a = append(a, strconv.Itoa(int(p.id)))             //id
	a = append(a, strconv.Itoa(p.sex))                 //sex
	a = append(a, toTime(p.dob).Format(dateLayoutISO)) //dob
	a = append(a, strconv.Itoa(p.age))                 //age
	a = append(a, toTime(p.regisDate).Format(dateLayoutISO))
	a = append(a, toTime(p.cancelDate).Format(dateLayoutISO))
	return a
}

func (p *Person) addVisits() {
	for _, disease := range p.config.Diseases {
		hadIt := p.sex == 0 && rand.Float64() < disease.PrevalenceMale ||
			p.sex == 1 && rand.Float64() < disease.PrevalenceFemale
		if !hadIt {
			continue
		}
		incidenceDate := RangeDate(p.regisDate, p.cancelDate)
		fup := (p.cancelDate - incidenceDate) / secondsInDay / daysInYear
		// estimate # of hospitalizations
		n := int64(Normal(disease.HospitalRate.Mean, disease.HospitalRate.SD)) * fup
		for i := int64(0); i < n; i++ {
			// p.visits = append(p.visits, p.newVisit(hospital, disease))
			p.dispatcher.SaveHosp(p.newVisit(hospital, disease).toStrings())
		}
		// estimate # of clinic encounters
		n = int64(Normal(disease.ClinicRate.Mean, disease.ClinicRate.SD)) * fup
		for i := int64(0); i < n; i++ {
			// p.visits = append(p.visits, p.newVisit(clinic, disease))
			p.dispatcher.SaveClinic(p.newVisit(clinic, disease).toStrings())
		}
	}
}

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

func (p *Person) newVisit(kind int, disease *Disease) *Visit {
	v := Visit{
		id:        p.id,
		startDate: RangeDate(p.regisDate, p.cancelDate),
	}
	switch kind {
	case hospital:
		v.endDate = v.startDate + int64(Normal(hospitalStayMean, hospitalStaySD))*secondsInDay
		v.diagnosis = disease.Icd10
	default:
		v.diagnosis = disease.Icd9
	}
	return &v
}
