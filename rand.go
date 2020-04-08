package main

import (
	"errors"
	"math/rand"
	"time"
)

// RangeInt returns an int in a range of two ints.
// it panics if max-min <0
func RangeInt(min, max int) int {
	// if max-min <= 0 {
	// 	log.Printf("RangeInt max %d min %d max-min %d", max, min, max-min)
	// }
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
	return rand.Int63n(max-min+1) + min
}

func toTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// func simNormal() {
// 	// Create a normal distribution
// 	dist := distuv.Normal{
// 		Mu:    2,
// 		Sigma: 5,
// 	}

// 	data := make([]float64, 1e5)

// 	// Draw some random values from the standard normal distribution
// 	for i := range data {
// 		data[i] = dist.Rand()
// 	}
// 	mean, std := stat.MeanStdDev(data, nil)
// 	meanErr := stat.StdErr(std, float64(len(data)))

// 	fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)
// }

// type Weibull struct {
// 	// Shape parameter of the distribution. A value of 1 represents
// 	// the exponential distribution. A value of 2 represents the
// 	// Rayleigh distribution. Valid range is (0,+∞).
// 	K float64
// 	// Scale parameter of the distribution. Valid range is (0,+∞).
// 	Lambda float64
// }

// // Rand returns a random sample drawn from the distribution.
// func (w Weibull) Rand() float64 {
// 	return w.Lambda * math.Pow(-math.Log(1-rand.Float64()), 1/w.K)
// }

// func genWeibull(k, lambda float64, size int64) []float64 {
// 	dist := Weibull{
// 		K:      k,
// 		Lambda: lambda,
// 	}

// 	data := make([]float64, size)
// 	// Draw some random values from the standard normal distribution
// 	for i := int64(0); i < size; i++ {
// 		data[i] = dist.Rand()
// 	}
// 	return data
// }

// func simWeibull(data []float64) {
// 	mean, std := stat.MeanStdDev(data, nil)
// 	meanErr := stat.StdErr(std, float64(len(data)))

// 	fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)
// }

// A Choice contains a generic item and a weight controlling the frequency with
// which it will be selected.
type Choice struct {
	Weight int
	Item   interface{}
}

// WeightedChoice used weighted random selection to return one of the supplied
// choices.  Weights of 0 are never selected.  All other weight values are
// relative.  E.g. if you have two choices both weighted 3, they will be
// returned equally often; and each will be returned 3 times as often as a
// choice weighted 1.
func WeightedChoice(choices []Choice) (Choice, error) {
	// Based on this algorithm:
	//     http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
	var ret Choice
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	r := RangeInt(0, sum)
	for _, c := range choices {
		r -= c.Weight
		if r < 0 {
			return c, nil
		}
	}
	err := errors.New("Internal error - code should not reach this point")
	return ret, err
}
