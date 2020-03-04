package simula

import (
	"math"
	"math/rand"
)

//TODO: use new source to avoid locking

func weibullRand(k, lambda float64) float64 {
	return lambda * math.Pow(-math.Log(1-rand.Float64()), 1/k)
}

// WeibullVector returns a random sample drawn from the distribution.
// k: Shape parameter of the distribution. A value of 1 represents
// the exponential distribution. A value of 2 represents the
// Rayleigh distribution. Valid range is (0,+∞).
// lambda Scale parameter of the distribution. Valid range is (0,+∞).
func WeibullVector(k, lambda float64, size int64) []float64 {
	data := make([]float64, size)
	for i := int64(0); i < size; i++ {
		data[i] = weibullRand(k, lambda)
	}
	return data
}


// WeibullVectorInt returns a random sample drawn from the distribution.
// k: Shape parameter of the distribution. A value of 1 represents
// the exponential distribution. A value of 2 represents the
// Rayleigh distribution. Valid range is (0,+∞).
// lambda Scale parameter of the distribution. Valid range is (0,+∞).
func WeibullVectorInt(k, lambda float64, size int64) []int {
	func weibullRandInt() int {
		return lambda * math.Pow(-math.Log(1-rand.Int()), 1/k)
	}
	data := make([]int, size)
	for i := int64(0); i < size; i++ {
		data[i] = weibullRandInt(k, lambda)
	}
	return data
}
