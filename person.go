package main

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	dateLayoutISO     = "2006-01-02"
	secondsInDay      = 24 * 60 * 60
	daysInYear        = 365
	stataMissingInt64 = 0x7fffffe5
)

const (
	kindHospital = iota
	kindClinic
)

var (
	today     = time.Now()
	todayUnix = today.Unix()
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
		dob:        RangeDate(config.Population.minDate, todayUnix),
		visits:     []*Visit{},
	}
	defer dispatcher.wg.Done()
	dob := toTime(p.dob)
	p.age = today.Year() - dob.Year()
	if rand.Float64() < config.Population.MigrantProb {
		p.regisDate = RangeDate(config.Population.databaseStartDate, todayUnix)
	} else {
		p.regisDate = config.Population.databaseStartDate
	}
	if p.dob > p.regisDate {
		p.dob = p.regisDate
	}
	if rand.Float64() < config.Population.CancelProb {
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
			p.dispatcher.SaveHosp(p.newVisit(kindHospital, disease, incidenceDate).toStrings())
		}
		// estimate # of clinic encounters
		n = int64(Normal(disease.ClinicRate.Mean, disease.ClinicRate.SD)) * fup
		for i := int64(0); i < n; i++ {
			p.dispatcher.SaveClinic(p.newVisit(kindClinic, disease, incidenceDate).toStrings())
		}
		// estimate # of Rxs filled
		n = int64(Normal(disease.RxRate.Mean, disease.RxRate.SD)) * fup
		for i := int64(0); i < n; i++ {
			p.dispatcher.SaveRx(p.newRx(disease, incidenceDate).toStrings())
		}
	}
}
