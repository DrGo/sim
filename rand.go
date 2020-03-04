package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func RangeInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// Normal returns a draw from a normally distributed dis with desired mean and sd
func Normal(mean, sd float64) float64 {
	return rand.NormFloat64()*sd + mean
}

// DateFromYear returns a valid date from a year and random month and day
// func DateFromYear(year int) time.Time {
// 	return rand.Intn(max-min+1) + min
// }

func RangeDate(min, max int64) int64 {
	return rand.Int63n((max - min)) + min
}

func toTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func simNormal() {
	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    2,
		Sigma: 5,
	}

	data := make([]float64, 1e5)

	// Draw some random values from the standard normal distribution
	for i := range data {
		data[i] = dist.Rand()
	}
	mean, std := stat.MeanStdDev(data, nil)
	meanErr := stat.StdErr(std, float64(len(data)))

	fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)
}

type Weibull struct {
	// Shape parameter of the distribution. A value of 1 represents
	// the exponential distribution. A value of 2 represents the
	// Rayleigh distribution. Valid range is (0,+∞).
	K float64
	// Scale parameter of the distribution. Valid range is (0,+∞).
	Lambda float64
}

// Rand returns a random sample drawn from the distribution.
func (w Weibull) Rand() float64 {
	return w.Lambda * math.Pow(-math.Log(1-rand.Float64()), 1/w.K)
}

func genWeibull(k, lambda float64, size int64) []float64 {
	dist := Weibull{
		K:      k,
		Lambda: lambda,
	}

	data := make([]float64, size)
	// Draw some random values from the standard normal distribution
	for i := int64(0); i < size; i++ {
		data[i] = dist.Rand()
	}
	return data
}

func simWeibull(data []float64) {
	mean, std := stat.MeanStdDev(data, nil)
	meanErr := stat.StdErr(std, float64(len(data)))

	fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)
}
