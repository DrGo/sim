package simula

import "math/rand"

// Returns an int >= min, < max
func UniformRangeRand(min, max int) int {
	return min + rand.Intn(max-min)
}
